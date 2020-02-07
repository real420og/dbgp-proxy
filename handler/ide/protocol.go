package idehandler

import (
	"encoding/xml"
	"fmt"
	"strings"
)

const commandInit = "proxyinit"
const commandStop = "proxystop"

type IdeCommand struct {
	Name   string
	Idekey string
	Port   string
}

type ProxyXmlMessage struct {
	XMLName xml.Name
	Success int8   `xml:"success,attr"`
	Idekey  string `xml:"idekey,attr"`
	Error   string `xml:"error>message"`
}

func newXmlMessage(ideCommand *IdeCommand) *ProxyXmlMessage {
	return &ProxyXmlMessage{
		XMLName: xml.Name{Local: ideCommand.Name},
		Success: 1,
		Idekey:  ideCommand.Idekey,
		Error:   "",
	}
}

func (xmlMessage *ProxyXmlMessage) setError(err error) {
	xmlMessage.Success = 0
	xmlMessage.Error = err.Error()
}

func createIdeCommand(commandString string) (*IdeCommand, error) {
	ideCommand := &IdeCommand{}

	pars := strings.Split(commandString, " ")

	ideCommand.Name = pars[0]
	if ideCommand.Name != commandInit && ideCommand.Name != commandStop {
		return ideCommand, fmt.Errorf("empty command name %s", ideCommand.Name)
	}

	l := len(pars)
	for i := 1; i < l; i++ {
		if pars[i] == "-p" {
			i++
			ideCommand.Port = pars[i]
		}

		if pars[i] == "-k" {
			i++
			ideCommand.Idekey = pars[i]
		}
	}

	return ideCommand, nil
}
