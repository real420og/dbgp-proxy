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

func (that *DebugHandler) Handle(conn net.Conn, xx *server.Xx) error {
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

	remoteAddr := conn.RemoteAddr().String()
	localAddr := conn.LocalAddr().String()

	//xx.Act6(remoteAddr)
	xx.Act5(remoteAddr+" through "+localAddr+" send "+string(message[:len(message)-1]))
	//xx.Act7(string(message[:len(message)-1]))

	idekey, err := getIdekey(message[:len(message)-1])
	key := "10.1.0.610.0.4.71:9002" + idekey

	d := &server.SerIdeKey{Server: conn.RemoteAddr().String(), IdeKey: idekey}
	xx.IdeWySendMessage(conn.RemoteAddr().String()+idekey, *d)

	if err != nil {
		return fmt.Errorf("%s", err)
	}

	log.Printf("idekey: %s", key)

	ideStorage, ok := that.listIdeConnection.FindIdeConnection(key)
	if !ok {
		return fmt.Errorf("client with idekey %s is not registered", key)
	}
	client, err := net.Dial("tcp", ideStorage.FullAddress())
	defer that.closer(client)

	if err != nil {
		return err
	}

	log.Println("IDE Connected")

	initMessage := append(dataLength, message...)
	_, err = io.Copy(client, bytes.NewReader(initMessage))
	if err != nil {
		return err
	}
	log.Println("init accepted")
	clientChan := make(chan error)
	serverChan := make(chan error)

	go func() {
		//xx.Act8(fmt.Sprintf("keep copy from %s to %s", conn.RemoteAddr().String(), conn.LocalAddr().String()))

		xx.Act8(fmt.Sprintf("keep copy from %s to %s", conn.RemoteAddr().String(), client.RemoteAddr().String()))

		//_, err := io.Copy(&b, io.TeeReader(client, conn))
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
			xx.Act8("clientChan err")
			return nil
		case err = <-serverChan:
			xx.Act8("serverChan  err")
			return nil
		default:
			time.Sleep(sleepTimeout)
		}
	}

}

func (that *DebugHandler) closer (closer io.Closer) {
	if closer == nil {
		return
	}

	err := closer.Close()
	if err != nil {
		log.Fatal(err)
	}
}