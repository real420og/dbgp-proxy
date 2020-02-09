package storage

import (
	"sync"
)

type ListIdeConnection struct {
	sync.Mutex
	ideConnection map[string]*IdeConnection
}

func NewListIdeConnection() *ListIdeConnection {
	return &ListIdeConnection{ideConnection: map[string]*IdeConnection{}}
}

func (list *ListIdeConnection) AddIdeConnection(ideConnection *IdeConnection) {
	list.Lock()
	list.ideConnection[ideConnection.key] = ideConnection
	list.Unlock()
}

func (list *ListIdeConnection) DeleteIdeConnection(key string) {
	list.Lock()
	delete(list.ideConnection, key)
	list.Unlock()
}

func (list *ListIdeConnection) FindIdeConnection(key string) (*IdeConnection, bool) {
	list.Lock()
	ideConnection, ok := list.ideConnection[key]
	list.Unlock()
	return ideConnection, ok
}

func (list *ListIdeConnection) HasIdeConnection(key string) bool {
	if _, ok := list.ideConnection[key]; ok {
		return true
	}

	return false
}

func (list *ListIdeConnection) HasNotIdeConnection(key string) bool {
	if _, ok := list.ideConnection[key]; ok {
		return false
	}

	return true
}
