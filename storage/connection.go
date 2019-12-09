package storage

import (
	"net"
)

type IdeConnection struct {
	address string
	port    string
	key     string
}

func NewIdeConnection(address string, port string, key string) *IdeConnection {
	return &IdeConnection{key: key, address: address, port: port}
}

func (ideConnection *IdeConnection) FullAddress() string {
	return net.JoinHostPort(ideConnection.address, ideConnection.port)
}
