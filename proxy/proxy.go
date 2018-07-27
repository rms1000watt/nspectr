package proxy

import (
	"fmt"
	"io"
	"net"
	"sync"

	log "github.com/sirupsen/logrus"
)

const (
	to = "127.0.0.1:8080"
)

// Proxy starts the proxy
func Proxy(cfg Config) {
	addr := fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)

	listener, err := net.Listen("tcp", addr)
	if err != nil {
		log.Error("Failed creating listener: ", err)
		return
	}

	doneCh := make(chan bool)
	go run(listener, doneCh)

	log.Info("Listening on: ", addr)
	<-doneCh
}

func run(listener net.Listener, doneCh chan bool) {
	for {
		select {
		case <-doneCh:
			return
		default:
			connection, err := listener.Accept()
			if err != nil {
				log.Error("Failed accepting connection: ", err)
				continue
			}

			go handle(connection, doneCh)
		}
	}
}

func handle(connection net.Conn, doneCh chan bool) {
	log.Debug("Start: Handling: ", connection)
	defer log.Debug("Done: Handling: ", connection)

	defer connection.Close()
	remote, err := net.Dial("tcp", to)
	if err != nil {
		log.Error("Failed dialing: ", err)
		return
	}
	defer remote.Close()

	wg := &sync.WaitGroup{}
	wg.Add(2)
	go copy(remote, connection, wg, doneCh)
	go copy(connection, remote, wg, doneCh)
	wg.Wait()
}

func copy(from, to net.Conn, wg *sync.WaitGroup, doneCh chan bool) {
	defer wg.Done()
	select {
	case <-doneCh:
		return
	default:
		if _, err := io.Copy(to, from); err != nil {
			log.Error("io.Copy failed: ", err)
			stop(doneCh)
			return
		}
	}
}

func stop(doneCh chan bool) {
	if doneCh == nil {
		return
	}
	close(doneCh)
	doneCh = nil
}
