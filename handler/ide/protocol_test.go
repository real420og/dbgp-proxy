package idehandler

import (
	"fmt"
	"reflect"
	"testing"
)

func Test_TestCommand(t *testing.T) {
	command := "test"
	ideCommand, err := createIdeCommand(command)

	if _, ok := err.(error); !ok {
		t.Errorf("expect err %s", command)
	}

	ideCommandMock := &IdeCommand{
		Name: command,
	}

	if !reflect.DeepEqual(ideCommand, ideCommandMock) {
		t.Errorf("value: %s expect %s", ideCommand, ideCommandMock)
	}
}

func Test_InitCommand(t *testing.T) {
	command := "proxyinit -p xxxx -k 8000"
	ideCommand, err := createIdeCommand(command)

	if _, ok := err.(error); ok {
		t.Errorf("expect err %s", command)
	}

	fmt.Printf("%s", ideCommand)

	ideCommandMock := &IdeCommand{
		Name:   "proxyinit",
		Idekey: "8000",
		Port:   "xxxx",
	}

	if !reflect.DeepEqual(ideCommand, ideCommandMock) {
		t.Errorf("value: %s expect %s", ideCommand, ideCommandMock)
	}
}
