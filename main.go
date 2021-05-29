package main

import (
	"flag"
	debughandler "github.com/real420og/dbgp-proxy/handler/debug"
	idehandler "github.com/real420og/dbgp-proxy/handler/ide"
	"github.com/real420og/dbgp-proxy/server"
	"github.com/real420og/dbgp-proxy/storage"
	"log"
	"net"
	"os"
	"os/signal"
	"sync"
)



func main() {
	ideHostnamePort := *flag.String("ide", "0.0.0.0:9002", "ip:port for ide connections")
	debugHostnamePort := *flag.String("xdebug", "0.0.0.0:9005", "ip:port for xdebug connections")
	flag.Parse()

	ideConnectionList := storage.NewListIdeConnection()

	syncGroup := &sync.WaitGroup{}


	c := &server.Xx{C: make(chan server.St)}

	ideServer := server.NewServer("ide", resolveTCP(ideHostnamePort), syncGroup, c)
	debugServer := server.NewServer("debug", resolveTCP(debugHostnamePort), syncGroup, c)

	go ideServer.Listen(idehandler.NewIdeHandler(ideConnectionList))
	go debugServer.Listen(debughandler.NewDebugHandler(ideConnectionList))

	go c.Read()
	//go func() {
	//	for i := range c.C {
	//		fmt.Print("\033[H\033[2J")
	//		if i.Act1 != nil {fmt.Println(i.Act1)} else{fmt.Println("")}
	//		if i.Act2 != nil {fmt.Println(i.Act2)} else{fmt.Println("")}
	//		if i.Act3 != nil {fmt.Println(i.Act3)} else{fmt.Println("")}
	//
	//		//fmt.Print("sssssssss: "+i)
	//	}
	//}()

	log.Println("dbgp proxy started")

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)

	log.Printf("signal: %s", <-signals)

	ideServer.Stop()
	debugServer.Stop()
	syncGroup.Wait()

	log.Println("dbgp proxy stopped")
}

func resolveTCP(host string) *net.TCPAddr {
	address, err := net.ResolveTCPAddr("tcp", host)
	if err != nil {
		log.Fatal(err)
	}
	return address
}
