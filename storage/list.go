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
	defer list.Unlock()
	list.ideConnection[ideConnection.key] = ideConnection
}

func (list *ListIdeConnection) DeleteIdeConnection(key string) {
	list.Lock()
	defer list.Unlock()
	delete(list.ideConnection, key)
}

func (list *ListIdeConnection) FindIdeConnection(key string) (*IdeConnection, bool) {
	list.Lock()
	defer list.Unlock()
	ideConnection, ok := list.ideConnection[key]
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
