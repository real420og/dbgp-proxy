package server

import (
	"io"
	"log"
	"net"
	"sync"
	"time"
)

const deadlineTime = time.Second * 2

type Server struct {
	name    string
	address *net.TCPAddr
	stop    bool
	group   *sync.WaitGroup
}

func NewServer(name string, address *net.TCPAddr, group *sync.WaitGroup) *Server {
	return &Server{
		name,
		address,
		false,
		group,
	}
}

func (server *Server) Stop() {
	server.stop = true
}

func (server *Server) Listen(handler Handler) {
	server.group.Add(1)
	defer server.group.Done()

	listener, err := net.ListenTCP("tcp", server.address)
	if err != nil {
		log.Fatal(err)
	}

	defer server.closeConnection(listener)

	log.Printf("start %s server on %s", server.name, server.address)

	for {
		if server.stop {
			break
		}
		_ = listener.SetDeadline(time.Now().Add(deadlineTime))
		conn, err := listener.Accept()
		if err != nil {
			if opErr, ok := err.(*net.OpError); ok && opErr.Timeout() {
				continue
			}
			log.Print(err)
			continue
		}
		go server.handleConnection(conn, handler)
	}

	log.Printf("shutdown %s server", server.name)
}

func (server *Server) handleConnection(conn net.Conn, handler Handler) {
	defer server.closeConnection(conn)
	server.group.Add(1)
	defer server.group.Done()
	log.Printf("start new connection on %s from %s", server.name, conn.RemoteAddr())
	err := handler.Handle(conn)

	if err != nil {
		log.Printf("handler response error %s", err)
	}

	log.Printf("close connection on %s from %s", server.name, conn.RemoteAddr())
}

func (server *Server) closeConnection(closer io.Closer) {
	err := closer.Close()
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("close %s connection", server.name)
}
