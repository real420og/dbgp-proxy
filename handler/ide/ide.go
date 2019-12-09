package idehandler

import (
	"bufio"
	"bytes"
	"encoding/xml"
	"fmt"
	"github.com/real420og/dbgp-proxy/server"
	"github.com/real420og/dbgp-proxy/storage"
	"io"
	"log"
	"net"
)

type IdeHandler struct {
	listIdeConnection *storage.ListIdeConnection
}

func NewIdeHandler(storage *storage.ListIdeConnection) *IdeHandler {
	return &IdeHandler{listIdeConnection: storage}
}

func (handler *IdeHandler) Handle(conn net.Conn) error {
	reader := bufio.NewReader(conn)
	data, err := reader.ReadBytes(server.ReadDelimiter)
	if err != nil && err != io.EOF {
		return fmt.Errorf("%s", err)
	}

	log.Println("new ide connection")

	message := string(data[:len(data)-1])

	command, err := createIdeCommand(message)
	xmlMessage := newXmlMessage(command)

	if err != nil {
		xmlMessage.setError(err)
		return handler.sendResponse(conn, xmlMessage)
	}

	host, _, err := net.SplitHostPort(conn.RemoteAddr().String())
	if err != nil {
		xmlMessage.setError(err)
		return handler.sendResponse(conn, xmlMessage)
	}

	err = handler.processMessage(command, host)

	if err != nil {
		xmlMessage.setError(err)
		return handler.sendResponse(conn, xmlMessage)
	}

	return handler.sendResponse(conn, xmlMessage)
}

func (handler *IdeHandler) processMessage(command *IdeCommand, host string) error {
	actions := map[string]Action{
		commandInit: &Init{
			host:    host,
			storage: handler.listIdeConnection,
		},
		commandStop: &Stop{
			storage: handler.listIdeConnection,
		},
	}

	if action, ok := actions[command.Name]; ok {
		err := action.Exec(command)

		if err, ok := err.(error); ok {
			return err
		}
	}

	return nil
}

func (handler *IdeHandler) sendResponse(conn net.Conn, proxyXmlMessage ProxyXmlMessage) error {
	message, err := xml.Marshal(proxyXmlMessage)

	if err != nil {
		return err
	}
	_, err = io.Copy(conn, bytes.NewReader(message))
	if err != nil {
		return err
	}
	return nil
}
