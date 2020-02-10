package idehandler

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

func Test_TestCommand(t *testing.T) {
	command := "test"
	ideCommand, err := createIdeCommand(command)
	assert.Errorf(t, err, "expect err %s", command)

	ideCommandMock := &IdeCommand{
		Name: command,
	}
	assert.Truef(t, reflect.DeepEqual(ideCommand, ideCommandMock), "value: %s expect %s", ideCommand, ideCommandMock)
}

func Test_InitCommand(t *testing.T) {
	command := "proxyinit -p xxxx -k 8000"
	ideCommand, err := createIdeCommand(command)
	assert.Errorf(t, err, "expect err %s", command)

	fmt.Printf("%s", ideCommand)
	ideCommandMock := &IdeCommand{
		Name:   "proxyinit",
		Idekey: "8000",
		Port:   "xxxx",
	}
	assert.Truef(t, reflect.DeepEqual(ideCommand, ideCommandMock), "value: %s expect %s", ideCommand, ideCommandMock)
}
