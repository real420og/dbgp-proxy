package server

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"strings"
	"sync"
	"time"
)

const deadlineTime = time.Second * 2

type Server struct {
	name    string
	address *net.TCPAddr
	stop    bool
	group   *sync.WaitGroup
	xx      *Xx
}

type Xx struct {
	C  chan St
	St St
}

func (xx *Xx) Act6(act string) {
	xx.St.Act6 = &act
	xx.C <- xx.St
}
func (xx *Xx) Act5(act string) {
	xx.St.Act5 = &act
	xx.C <- xx.St
}
func (xx *Xx) Act7(act string) {
	xx.St.Act7 = &act
	xx.C <- xx.St
}
func (xx *Xx) Act8(act string) {
	xx.St.Act8 = &act
	xx.C <- xx.St
}
func (xx *Xx) Act9(act string) {
	xx.St.Act9 = &act
	xx.C <- xx.St
}
func (xx *Xx) Act1(act1 string) {
	xx.St.Act1 = &act1
	xx.C <- xx.St
}
func (xx *Xx) Act2(act2 string) {
	xx.St.Act2 = &act2
	xx.C <- xx.St
}
func (xx *Xx) Act3(act3 string) {
	xx.St.Act3 = &act3
	xx.C <- xx.St
}
func (xx *Xx) Act4(act4 string) {
	xx.St.Act4 = &act4
	xx.C <- xx.St
}
func (xx *Xx) IdeWySendMessage(key string, act SerIdeKey) {
	xx.St.IdeWySendMessage[key] = &act
	xx.C <- xx.St
}
func (xx *Xx) Read() {
	for i := range xx.C {
		fmt.Print("\033[H\033[2J")
		if i.Act1 != nil {fmt.Println(fmt.Sprintf("1: %s \n\n", *i.Act1))} else{fmt.Println("")}
		if i.Act2 != nil {fmt.Println(fmt.Sprintf("2: %s \n\n", *i.Act2))} else{fmt.Println("")}
		if i.Act3 != nil {fmt.Println(fmt.Sprintf("3: %s \n\n", *i.Act3))} else{fmt.Println("")}
		if i.Act4 != nil {fmt.Println(fmt.Sprintf("4: %s \n\n", *i.Act4))} else{fmt.Println("")}
		if i.Act5 != nil {fmt.Println(fmt.Sprintf("5: %s \n\n", *i.Act5))} else{fmt.Println("")}
		if i.Act6 != nil {fmt.Println(fmt.Sprintf("6: %s \n\n", *i.Act6))} else{fmt.Println("")}
		if i.Act7 != nil {fmt.Println(fmt.Sprintf("7: %s \n\n", *i.Act7))} else{fmt.Println("")}
		if i.Act8 != nil {fmt.Println(fmt.Sprintf("8: %s \n\n", *i.Act8))} else{fmt.Println("")}
		if i.Act9 != nil {fmt.Println(fmt.Sprintf("9: %s \n\n", *i.Act9))} else{fmt.Println("")}

		if i.IdeWySendMessage != nil {
			dd, _ := json.Marshal(i.IdeWySendMessage)
			if i.IdeWySendMessage != nil {
				fmt.Println(fmt.Sprintf("IdeWySendMessage: %s \n\n", dd))
			} else {
				fmt.Println("")
			}
		}
	}
}

type St struct {
	Act1 *string
	Act2 *string
	Act3 *string
	Act4 *string
	Act6 *string
	Act5 *string
	Act7 *string
	Act8 *string
	Act9 *string
	IdeWySendMessage map[string]*SerIdeKey
}


type SerIdeKey struct {
	Server string
	IdeKey string
}

func NewServer(name string, address *net.TCPAddr, group *sync.WaitGroup, xx *Xx) *Server {
	return &Server{
		name,
		address,
		false,
		group,
		xx,
	}
}

func (that *Server) Stop() {
	that.stop = true
}

func (that *Server) Listen(handler Handler) {
	that.group.Add(1)
	defer that.group.Done()

	listener, err := net.ListenTCP("tcp", that.address)
	if err != nil {
		log.Fatal(err)
	}

	that.xx.Act1("text value 3245")


	defer that.closeConnection(listener)

	log.Printf("start %s server on %s", that.name, that.address)

	for {
		that.xx.Act1("loop")
		if that.stop {
			break
		}
		_ = listener.SetDeadline(time.Now().Add(deadlineTime))
		conn, err := listener.Accept()
		if err != nil {
			if opErr, ok := err.(*net.OpError); ok && opErr.Timeout() {

				a := []string{
					"accept tcp nil",
					opErr.Addr.Network(),
					opErr.Addr.String(),
					opErr.Error(),
					opErr.Net,
					opErr.Op,
				}

				that.xx.Act2(strings.Join(a, " "))
				continue
			}
			log.Print(err)
			continue
		}

		go func(conn net.Conn, handler Handler, xx *Xx) {
			defer that.closeConnection(conn)
			that.group.Add(1)
			defer that.group.Done()

			err := handler.Handle(conn, xx)

			if err != nil {
				xx.Act3(fmt.Sprintf("handler response error: %s", err))
			}
		}(conn, handler, that.xx)
	}

	log.Printf("shutdown %s server", that.name)
}

func (that *Server) closeConnection(closer io.Closer) {
	err := closer.Close()
	if err != nil {
		log.Fatal(err)
	}
}
