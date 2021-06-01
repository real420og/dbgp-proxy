package idehandler

import (
	"encoding/json"
	"fmt"
	"github.com/real420og/dbgp-proxy/server"
	"github.com/real420og/dbgp-proxy/storage"
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
	key := s.IdeHost + s.LocalAddr + command.Idekey
	if s.Storage.HasNotIdeConnection(key) {

		s.Storage.AddIdeConnection(key, &storage.IdeConnection{
			IdeHost: s.IdeHost,
			RemoteAddr: s.RemoteAddr,
			LocalAddr: s.LocalAddr,
			IdeKey: command.Idekey,
			Port: command.Port,
		})

		ss, _ := json.Marshal(s.Storage)
		xx.Act4(string(ss))

		return nil
	}

	return fmt.Errorf("already in use")
}

type Stop struct {
	storage    *storage.ListIdeConnection
	IdeHost    string
	RemoteAddr string
	LocalAddr  string
	Storage    *storage.ListIdeConnection
}

func (s *Stop) Exec(command *IdeCommand, xx *server.Xx) error {
	key := s.IdeHost + s.LocalAddr + command.Idekey
	if s.Storage.HasIdeConnection(key) {
		s.Storage.DeleteIdeConnection(key)

		ss, _ := json.Marshal(s.Storage)
		xx.Act4(string(ss))

		return nil
	}

	//log.Printf("attempt to delete idekey %s", command.Idekey)
	return fmt.Errorf("is not registered")
}
