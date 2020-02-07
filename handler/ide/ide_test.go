package idehandler

import (
	"encoding/xml"
	"github.com/real420og/dbgp-proxy/storage"
	"github.com/stretchr/testify/assert"
	"net"
	"testing"
)

func TestSendResponsef(t *testing.T) {
	ideConnectionList := storage.NewListIdeConnection()
	handler := NewIdeHandler(ideConnectionList)

	command, err := createIdeCommand("proxyinit")
	assert.NoError(t, err)
	pxm := &ProxyXmlMessage{
		XMLName: xml.Name{Local: command.Name},
		Success: 1,
		Idekey:  command.Idekey,
		Error:   "",
	}

	connx := &net.TCPConn{}
	err = handler.sendResponse(connx, pxm)
	assert.Errorf(t, err, "invalid argument")
}
