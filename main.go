package main

import (
	"flag"
	"log"
	"message/cmd"
	"message/internal/util/util_net"
	"os"
	"runtime"
	"time"
)

//var _mode =

var (
	BUILD_DATE = "RandSGVsbG8gV29ypm"
	GIT_BRANCH = "RandBGQg2xzY2FXzGK"
	GIT_COMMIT = "RandGZFka2xzGFkc2E"
)

func init() {
	log.Println("main [init] start ")
	defer func() {
		log.Println("main [init] end ")
	}()

	for true {
		bAllOK, _ := util_net.CheckTcpService(
			[]string{
				"redis-service:6379", // todo: coding style
				"rmq-message:5672",
			},
		)

		log.Println("bAllOK: ", bAllOK)
		if bAllOK {
			break
		}

		time.Sleep(time.Second)
	}
}

func main() {
	modeStr := ""
	defer func() {
		if r := recover(); r != nil {
			log.Println("Recovered:", r)
			stack := make([]byte, 1024*4)
			runtime.Stack(stack, false)
			log.Printf("stack: %s\n", stack)
		}
	}()

	// debug
	log.Println("BUILD_DATE: ", BUILD_DATE)
	log.Println("GIT_BRANCH: ", GIT_BRANCH)
	log.Println("GIT_COMMIT: ", GIT_COMMIT)

	value := os.Getenv("MODE")
	if value == "" { // config.yaml
		flag.Parse()
		modeStr = *flag.String("mode", "empty", "work mode")
	} else { // k8s
		modeStr = value
	}

	log.Println("modeStr: ", modeStr)

	switch modeStr {
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
