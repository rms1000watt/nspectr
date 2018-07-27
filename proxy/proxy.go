package proxy

import (
	"fmt"
	"io"
	"net"

	log "github.com/sirupsen/logrus"
)

const (
	to = "127.0.0.1:8080"
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

		go handle(client)
	}
}

func handle(client *net.TCPConn) {
	log.Debugf("Start: %s->%s", client.RemoteAddr(), to)
	defer log.Debugf("Done: %s->%s", client.RemoteAddr(), to)

	server, err := getServer()
	if err != nil {
		log.Error("Failed getting server: ", err)
		return
	}

	serverDone := make(chan bool)
	clientDone := make(chan bool)

	go copy(server, client, clientDone)
	go copy(client, server, serverDone)

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

func copy(dst, src net.Conn, done chan bool) {
	if _, err := io.Copy(dst, src); err != nil {
		log.Error("io.Copy(dst, src) failed: ", err)
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

func getServer() (server *net.TCPConn, err error) {
	serverAddr, err := net.ResolveTCPAddr("tcp", to)
	if err != nil {
		log.Debug("Failed resolving tcp addr: ", err)
		return
	}

	return net.DialTCP("tcp", nil, serverAddr)
}
