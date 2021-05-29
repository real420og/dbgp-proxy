package storage

import (
	"sync"
)

type ListIdeConnection struct {
	sync.Mutex
	IdeConnection map[string]*IdeConnection
}

func NewListIdeConnection() *ListIdeConnection {
	return &ListIdeConnection{IdeConnection: map[string]*IdeConnection{}}
}

func (list *ListIdeConnection) AddIdeConnection(key string, ideConnection *IdeConnection) {
	list.Lock()
	list.IdeConnection[key] = ideConnection
	list.Unlock()
}

func (list *ListIdeConnection) DeleteIdeConnection(key string) {
	list.Lock()
	delete(list.IdeConnection, key)
	list.Unlock()
}

func (list *ListIdeConnection) FindIdeConnection(key string) (*IdeConnection, bool) {
	list.Lock()
	ideConnection, ok := list.IdeConnection[key]
	list.Unlock()
	return ideConnection, ok
}

func (list *ListIdeConnection) HasIdeConnection(key string) bool {
	if _, ok := list.IdeConnection[key]; ok {
		return true
	}

	return false
}

func (list *ListIdeConnection) HasNotIdeConnection(key string) bool {
	if _, ok := list.IdeConnection[key]; ok {
		return false
	}

	return true
}
