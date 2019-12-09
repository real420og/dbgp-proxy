package server

import (
	"net"
)

const ReadDelimiter = 0

type Handler interface {
	Handle(conn net.Conn) error
}
