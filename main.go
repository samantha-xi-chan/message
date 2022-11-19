package main

import (
	"flag"
	"github.com/message/cmd"
	"log"
	"net/http"
)

var _mode = flag.String("mode", "empty", "work mode")

const (
	BUILD_DATE = "RandSGVsbG8gV29ypm"
	GIT_BRANCH = "RandBGQg2xzY2FXzGK"
	GIT_COMMIT = "RandGZFka2xzGFkc2E"
)

func main() {
	// debug
	log.Println("BUILD_DATE: ", BUILD_DATE)
	log.Println("GIT_BRANCH: ", GIT_BRANCH)
	log.Println("GIT_COMMIT: ", GIT_COMMIT)

	go func() {
		http.ListenAndServe("0.0.0.0:9999", nil)
	}()

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
