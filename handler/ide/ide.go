package idehandler

import (
	"bufio"
	"bytes"
	"encoding/xml"
	"fmt"
	"github.com/real420og/dbgp-proxy/server"
	"github.com/real420og/dbgp-proxy/storage"
	"io"
	"net"
	"strings"
)

type IdeHandler struct {
	listIdeConnection *storage.ListIdeConnection
}

func NewIdeHandler(storage *storage.ListIdeConnection) *IdeHandler {
	return &IdeHandler{listIdeConnection: storage}
}

func (that *IdeHandler) Handle(conn net.Conn, xx *server.Xx) error {
	reader := bufio.NewReader(conn)
	data, err := reader.ReadBytes(server.ReadDelimiter)
	if err != nil && err != io.EOF {
		return fmt.Errorf("%s", err)
	}

	xx.Act3(strings.Join([]string{"ide handler: ", fmt.Sprintf("%s", data)}, " "))

	command, err := createIdeCommand(string(data[:len(data)-1]))
	xmlMessage := newXmlMessage(command)

	if err != nil {
		return that.sendResponse(conn, xmlMessage, err)
	}

	remoteAddr := conn.RemoteAddr().String()
	localAddr := conn.LocalAddr().String()

	ideHost, _, err := net.SplitHostPort(conn.RemoteAddr().String())
	if err != nil {
		return that.sendResponse(conn, xmlMessage, err)
	}

	actions := map[string]Execer{
		commandInit: &Init{
			IdeHost:    ideHost,
			RemoteAddr:    remoteAddr,
			LocalAddr:    localAddr,
			Storage: that.listIdeConnection,
		},
		commandStop: &Stop{
			IdeHost:    ideHost,
			RemoteAddr:    remoteAddr,
			LocalAddr:    localAddr,
			Storage: that.listIdeConnection,
		},
	}

	err = that.processMessage(command, actions, xx)

	return that.sendResponse(conn, xmlMessage, err)
}

func (that *IdeHandler) processMessage(command *IdeCommand, actions map[string]Execer, xx  *server.Xx) error {
	if action, ok := actions[command.Name]; ok {
		err := action.Exec(command, xx)

		if err, ok := err.(error); ok {
			return err
		}
	}

	return nil
}

func (that *IdeHandler) sendResponse(conn net.Conn, proxyXmlMessage *ProxyXmlMessage, err error) error {
	if err != nil {
		proxyXmlMessage.setError(err)
	}

	message, err := xml.Marshal(&proxyXmlMessage)

	if err != nil {
		return err
	}
	_, err = io.Copy(conn, bytes.NewReader(message))
	if err != nil {
		return err
	}
	return nil
}
