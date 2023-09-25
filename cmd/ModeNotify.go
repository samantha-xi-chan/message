package cmd

import (
	"encoding/json"
	mapset "github.com/deckarep/golang-set"
	"github.com/gorilla/websocket"
	"github.com/streadway/amqp"
	"log"
	api "message/api/domain"
	"message/internal/domain"
	"message/internal/service"
	"message/package/util_debug"

	"message/internal/config"
	"net/http"
)

const (
//
)

var (
	WildcardAsterisk = "*"
	mapTopicChanSet  = make(map[string](mapset.Set))
)

var upgrade = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func socketHandlerB(w http.ResponseWriter, r *http.Request) { // block if connection work
	var handlerTopic = "default"
	conn, err := upgrade.Upgrade(w, r, nil)
	if err != nil {
		log.Print("Error during connection upgrade:", err)
		return
	}
	defer conn.Close()

	wMsgChan := make(chan []byte, config.BUF_SIZE)
	wCtrlChan := make(chan []byte, config.BUF_SIZE)

	go func() {
		for {
			select {
			case msg, ok := <-wCtrlChan:
				if ok {
					err = conn.WriteMessage(websocket.TextMessage, msg)
					if err != nil {
						log.Println("Error during message writing: wCtrlChan", err)
						break
					}
				}
			case msg, ok := <-wMsgChan:
				if ok {
					err = conn.WriteMessage(websocket.TextMessage, msg)
					if err != nil {
						log.Println("Error during message writing wMsgChan:", err)
						break
					}
				}
			} // end select
		} // end for
	}()

	for {
		_, payload, err := conn.ReadMessage()
		if err != nil {
			log.Println("Error during message reading:", err)
			break
		}
		log.Printf("Received: %s", payload) // debug level

		if string(payload) == "ping" {
			wCtrlChan <- payload
			continue
		}

		req := api.WsReq{}
		e := json.Unmarshal(payload, &req)
		if e != nil {
			log.Println("Unmarshal err: ", e)
			continue
		}

		if req.Type == config.TYPE_SUBSCRIBE && req.Version == config.WS_PROTO_VER {
			sub := &api.Subscribe{}
			wsReq := api.WsReq{Payload: sub}
			json.Unmarshal(payload, &wsReq)

			currentTopic := sub.Topic
			log.Printf("handlerTopic:%s topic: %s", handlerTopic, currentTopic)

			if currentTopic == WildcardAsterisk {
				mapTopicChanSet[WildcardAsterisk].Add(wMsgChan)
			} else {
				if mapTopicChanSet[currentTopic] == nil {
					mapTopicChanSet[currentTopic] = mapset.NewSet()
				}
				mapTopicChanSet[currentTopic].Add(wMsgChan)
			}
			handlerTopic = currentTopic
		} else {
			log.Println("req.Type and ver not supported: ")
			continue
		}
	}

	mapTopicChanSet[handlerTopic].Remove(wMsgChan)

	//fmt.Println("socketHandlerB ending...")
}

func processMessages(ch <-chan amqp.Delivery, label string) {
	for d := range ch {
		log.Printf("%s: %s", label, d.Body)
		bytes := []byte(d.Body)
		broadcast(bytes)
	}
}

func MainModeNotify() {
	debugOn, e := config.GetDebugMode()
	if e != nil {
		log.Fatal("GetDebugMode: ", e)
	}

	if debugOn {
		addr, e := config.GetDebugPprofNotify()
		if e != nil {
			log.Fatal("config e: ", e)
		}

		log.Println("GetDebugPprofNotify addr: ", addr)
		go util_debug.InitPProf(addr)
	}

	mapTopicChanSet[WildcardAsterisk] = mapset.NewSet()

	v, _ := config.GetDependQueue()
	msgHigh, _ := service.QueueConnInit(v, config.EXCHANGE_HIGH)
	msgNormal, _ := service.QueueConnInit(v, config.EXCHANGE_NORMAL)

	//go func() {
	//	for d := range msgHigh {
	//		log.Printf("msgHigh.d: %s", d.Body)
	//
	//		bytes := []byte(d.Body)
	//		broadcast(bytes)
	//	}
	//}()
	//go func() {
	//	for d := range msgNormal {
	//		log.Printf("msgNormal.new=%s", d.Body)
	//		bytes := []byte(d.Body)
	//		broadcast(bytes)
	//	}
	//}()

	go func() {
		processMessages(msgHigh, "msgHigh")
	}()

	go func() {
		processMessages(msgNormal, "msgNormal")
	}()

	http.HandleFunc("/api/v1/socket", socketHandlerB)

	val, _ := config.GetNotifyPortHttp()
	log.Fatal(http.ListenAndServe(val, nil))
}

func broadcast(bytes []byte) {
	var feedSessionStream domain.FeedSessionStream
	err := json.Unmarshal(bytes, &feedSessionStream)
	if err != nil {
		log.Println("Error", err)
	} else {
		log.Println("feedSessionStream: ", feedSessionStream)
	}
	sessionId := feedSessionStream.SessionID
	if mapTopicChanSet[sessionId] != nil {
		for val := range mapTopicChanSet[sessionId].Iterator().C {
			val.(chan []byte) <- bytes
		}
	}
	// if WILDCARD_ASTERISK enabled
	for val := range mapTopicChanSet[WildcardAsterisk].Iterator().C {
		val.(chan []byte) <- bytes
	}

}
