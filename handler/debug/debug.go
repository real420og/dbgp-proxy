package debughandler

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/real420og/dbgp-proxy/server"
	"github.com/real420og/dbgp-proxy/storage"
	"io"
	"log"
	"net"
	"strconv"
	"time"
)

const sleepTimeout = time.Millisecond * 50

type DebugHandler struct {
	listIdeConnection *storage.ListIdeConnection
}

func NewDebugHandler(storage *storage.ListIdeConnection) *DebugHandler {
	return &DebugHandler{listIdeConnection: storage}
}

func (proxy *DebugHandler) Handle(conn net.Conn, xx *server.Xx) error {
	reader := bufio.NewReader(conn)
	dataLength, err := reader.ReadBytes(server.ReadDelimiter)
	if err != nil && err != io.EOF {
		return fmt.Errorf("%s", err)
	}

	if _, err := strconv.Atoi(string(dataLength[:len(dataLength)-1])); err != nil {
		return fmt.Errorf("%s", err)
	}

	message, err := reader.ReadBytes(server.ReadDelimiter)
	if err != nil && err != io.EOF {
		return fmt.Errorf("%s", err)
	}

	xx.Act1(string(message[:len(message)-1]))

	idekey, err := getIdekey(message[:len(message)-1])
	if err != nil {
		return fmt.Errorf("%s", err)
	}

	log.Printf("idekey: %s", idekey)
	err = proxy.sendAndPipe(conn, idekey, append(dataLength, message...))
	if err != nil {
		log.Println(err)
	}

	return nil
}

func (proxy *DebugHandler) sendAndPipe(conn net.Conn, idekey string, initMessage []byte) error {
	ideConnection, ok := proxy.listIdeConnection.FindIdeConnection(idekey)
	if !ok {
		return fmt.Errorf("client with idekey %s is not registered", idekey)
	}
	log.Printf("send init packet to: %s", ideConnection.FullAddress())
	client, err := net.Dial("tcp", ideConnection.FullAddress())

	defer func(closer io.Closer) {
		if closer == nil {
			return
		}

		err := closer.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(client)

	if err != nil {
		return err
	}

	log.Println("IDE Connected")
	_, err = io.Copy(client, bytes.NewReader(initMessage))
	if err != nil {
		return err
	}
	log.Println("init accepted")
	clientChan := make(chan error)
	serverChan := make(chan error)
	go func() {
		_, err = io.Copy(client, conn)
		clientChan <- err
	}()

	go func() {
		_, err = io.Copy(conn, client)
		serverChan <- err
	}()

	for {
		select {
		case err = <-clientChan:
			log.Println("IDE connection closed")
			log.Println("stop piping")
			return nil
		case err = <-serverChan:
			log.Println("XDebug connection closed")
			log.Println("stop piping")
			return nil
		default:
			time.Sleep(sleepTimeout)
		}
	}
}
