package main

import (
	"flag"
	"github.com/message/cmd"
	"log"
)

var _mode = flag.String("mode", "empty", "work mode")

func main() {
	flag.Parse()

	log.Println("current mode:", *_mode)

	switch *_mode {
	case "waiter":
		cmd.MainModeWaiter()
	case "notify":
		cmd.MainModeNotify()
	case "sink":
		cmd.MainModeSink()
	case "gateway":
		cmd.MainModeGateway()
	default:
		log.Println("usage: message -mode waiter/notify/sink/gateway")
	}

}
