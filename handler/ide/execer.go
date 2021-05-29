package idehandler

import (
	"fmt"
	"github.com/real420og/dbgp-proxy/server"
	"github.com/real420og/dbgp-proxy/storage"
	"log"
)

type Execer interface {
	Exec(command *IdeCommand, xx *server.Xx) error
}

type Init struct {
	IdeHost    string
	RemoteAddr string
	LocalAddr  string
	Storage    *storage.ListIdeConnection
}

func (s *Init) Exec(command *IdeCommand, xx *server.Xx) error {
	if s.Storage.HasNotIdeConnection(command.Idekey) {

		dd := &storage.IdeConnection{IdeHost: s.IdeHost, RemoteAddr: s.RemoteAddr, LocalAddr: s.LocalAddr, IdeKey: command.Idekey, Port: command.Port}
		s.Storage.AddIdeConnection(command.Idekey, dd)

		return nil
	}

	return fmt.Errorf("%s", command.Idekey)
}

type Stop struct {
	storage    *storage.ListIdeConnection
	IdeHost    string
	RemoteAddr string
	LocalAddr  string
	Storage    *storage.ListIdeConnection
}

func (s *Stop) Exec(command *IdeCommand, xx *server.Xx) error {
	if s.Storage.HasIdeConnection(command.Idekey) {
		s.Storage.DeleteIdeConnection(command.Idekey)
		log.Printf("delete client with idekey %s", command.Idekey)
		return nil
	}

	log.Printf("attempt to delete idekey %s", command.Idekey)
	return fmt.Errorf("idekey %s is not registered", command.Idekey)
}
