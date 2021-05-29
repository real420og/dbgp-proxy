package server

import (
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
func (xx *Xx) Read() {
	for i := range xx.C {
		fmt.Print("\033[H\033[2J")
		if i.Act1 != nil {fmt.Println(*i.Act1)} else{fmt.Println("")}
		if i.Act2 != nil {fmt.Println(*i.Act2)} else{fmt.Println("")}
		if i.Act3 != nil {fmt.Println(*i.Act3)} else{fmt.Println("")}
		if i.Act4 != nil {fmt.Println(*i.Act4)} else{fmt.Println("")}
	}
}

type St struct {
	Act1 *string
	Act2 *string
	Act3 *string
	Act4 *string
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
