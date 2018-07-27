package proxy

import (
	"bytes"
	"fmt"
	"io"
	"net"

	log "github.com/sirupsen/logrus"
)

const (
	killStr = "XXXXXX"
)

// Proxy starts the proxy
func Proxy(cfg Config) {
	addr := fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)

	listener, err := getListener(addr)
	if err != nil {
		log.Error("Failed getting listener: ", err)
		return
	}
	defer listener.Close()

	log.Info("Listening on: ", addr)
	for {
		client, err := listener.AcceptTCP()
		if err != nil {
			log.Error("Failed accepting client: ", err)
			continue
		}

		go handle(cfg, client)
	}
}

func handle(cfg Config, client *net.TCPConn) {
	log.Debugf("Start: %s->%s", client.RemoteAddr(), cfg.BackendAddr)
	defer log.Debugf("Done: %s->%s", client.RemoteAddr(), cfg.BackendAddr)

	server, err := getServer(cfg)
	if err != nil {
		log.Error("Failed getting server: ", err)
		return
	}

	serverDone := make(chan bool)
	clientDone := make(chan bool)

	go cpKill(server, client, clientDone)
	go cp(client, server, serverDone)

	var waitFor chan bool
	select {
	case <-clientDone:
		server.SetLinger(0)
		server.CloseRead()
		server.CloseWrite()
		waitFor = serverDone
	case <-serverDone:
		client.CloseRead()
		client.CloseWrite()
		waitFor = clientDone
	}

	<-waitFor
}

func cp(dst, src net.Conn, done chan bool) {
	b := make([]byte, 6)
	for {
		_, err := src.Read(b)
		if err == io.EOF {
			break
		}
		dst.Write(b)
	}

	if err := src.Close(); err != nil {
		log.Error("src.Close() failed: ", err)
	}

	done <- true
}

func cpKill(dst, src net.Conn, done chan bool) {
	b := make([]byte, 6)
	kill := make([]byte, 6)
	copy(kill[:], killStr)

	for {
		_, err := src.Read(b)
		if err == io.EOF {
			break
		}

		if bytes.Equal(b, kill) {
			log.Debug("kill sequence found")
			break
		}

		dst.Write(b)
	}

	if err := src.Close(); err != nil {
		log.Error("src.Close() failed: ", err)
	}

	done <- true
}

func getListener(addr string) (listener *net.TCPListener, err error) {
	listenAddr, err := net.ResolveTCPAddr("tcp", addr)
	if err != nil {
		log.Debug("Failed resolving tcp addr: ", err)
		return
	}

	return net.ListenTCP("tcp", listenAddr)
}

func getServer(cfg Config) (server *net.TCPConn, err error) {
	serverAddr, err := net.ResolveTCPAddr("tcp", cfg.BackendAddr)
	if err != nil {
		log.Debug("Failed resolving tcp addr: ", err)
		return
	}

	return net.DialTCP("tcp", nil, serverAddr)
}
