package debughandler

import (
	"bytes"
	"encoding/xml"
	"golang.org/x/net/html/charset"
)

type DebugPacket struct {
	Idekey string `xml:"idekey,attr"`
	Appid string `xml:"appid,attr"`
}

func getIdekey(packet []byte) (string, error) {
	var debugPacket DebugPacket
	decoder := xml.NewDecoder(bytes.NewReader(packet))
	decoder.CharsetReader = charset.NewReaderLabel
	err := decoder.Decode(&debugPacket)
	if err != nil {
		return "", err
	}
	return debugPacket.Idekey, nil
}
