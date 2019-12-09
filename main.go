package main

import (
	"flag"
	"github.com/real420og/dbgp-proxy/handler/debug"
	"github.com/real420og/dbgp-proxy/handler/ide"
	"github.com/real420og/dbgp-proxy/server"
	"github.com/real420og/dbgp-proxy/storage"
	"log"
	"net"
	"os"
	"os/signal"
	"sync"
)

func main() {
	ideHostnamePort := flag.String("ide", "0.0.0.0:9002", "ip:port for ide connections")
	debugHostnamePort := flag.String("xdebug", "0.0.0.0:9001", "ip:port for xdebug connections")
	flag.Parse()

	ideConnectionList := storage.NewListIdeConnection()

	syncGroup := &sync.WaitGroup{}
	ideServer := server.NewServer("ide", resolveTCP(*ideHostnamePort), syncGroup)
	debugServer := server.NewServer("debug", resolveTCP(*debugHostnamePort), syncGroup)

	go ideServer.Listen(idehandler.NewIdeHandler(ideConnectionList))
	go debugServer.Listen(debughandler.NewDebugHandler(ideConnectionList))

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
