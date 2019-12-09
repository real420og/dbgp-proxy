package idehandler

import (
	"fmt"
	"github.com/real420og/dbgp-proxy/storage"
	"log"
)

type Action interface {
	Exec(command *IdeCommand) error
}

type Init struct {
	host    string
	storage *storage.ListIdeConnection
}

func (s *Init) Exec(command *IdeCommand) error {
	if s.storage.HasNotIdeConnection(command.Idekey) {
		s.storage.AddIdeConnection(storage.NewIdeConnection(s.host, command.Port, command.Idekey))
		log.Printf("add client with host %s, port %s and idekey %s", s.host, command.Port, command.Idekey)
		return nil
	}

	return fmt.Errorf("ide %s already in use", command.Idekey)
}

type Stop struct {
	storage *storage.ListIdeConnection
}

func (s *Stop) Exec(command *IdeCommand) error {
	if s.storage.HasIdeConnection(command.Idekey) {
		s.storage.DeleteIdeConnection(command.Idekey)
		log.Printf("delete client with idekey %s", command.Idekey)
		return nil
	}

	log.Printf("attempt to delete idekey %s", command.Idekey)
	return fmt.Errorf("idekey %s is not registered", command.Idekey)
}
