package cmd

import (
	"encoding/json"
	"fmt"
	mapset "github.com/deckarep/golang-set"
	"github.com/gorilla/websocket"
	"log"
	api "message/api/domain"
	"message/internal/domain"
	"message/internal/service"
	"message/package/util_debug"

	"message/internal/config"
	"net/http"
)

var map_topic_chanset map[string](mapset.Set)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func socketHandlerB(w http.ResponseWriter, r *http.Request) { // block if connection work
	var topic = "default"
	//upgrader.CheckOrigin = func(r *http.Request) bool { return true }
	// Upgrade our raw HTTP connection to a websocket based one
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("Error during connection upgrade:", err)
		return
	}
	defer conn.Close()

	//r_chan := make(chan []byte, 2)
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
		log.Printf("Received: %s", payload)

		if string(payload) == "ping" {
			wCtrlChan <- payload
			continue
		}

		// 解析是否是订阅,如果是订阅则新增映射
		req := api.WsReq{}
		e := json.Unmarshal(payload, &req)
		if e != nil {
			log.Println("Unmarshal err: ", e)
			continue
		}

		if req.Type == config.TYPE_SUBSCRIBE && req.Version == config.WS_PROTO_VER {
			//sub := api.Subscribe{}

			// unmarshal it (usually after receiving bytes from somewhere)
			sub := &api.Subscribe{}
			wsReq := api.WsReq{Payload: sub}
			json.Unmarshal(payload, &wsReq)

			topic = sub.Topic
			log.Printf("topic: %s", topic)
			if map_topic_chanset[topic] == nil {
				log.Printf("new map KV: %s", topic)
				map_topic_chanset[topic] = mapset.NewSet()
			}
			map_topic_chanset[topic].Add(wMsgChan)

		} else {
			log.Println("req.Type and ver not supported: ")
			continue
		}
	}

	map_topic_chanset[topic].Remove(wMsgChan)

	fmt.Println("socketHandlerB ending...")
	select {}
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

	v, _ := config.GetDependQueue()

	msgs_high, _ := service.QueueConnInit(v, config.EXCHANGE_HIGH)
	msgs_normal, _ := service.QueueConnInit(v, config.EXCHANGE_NORMAL)

	map_topic_chanset = make(map[string](mapset.Set))
	//map_topic_chanset[TOPIC] = mapset.NewSet() // todo: 移除 改为 后期新增订阅时 如果没有这条KV映射则新增set并增加map映射， 如果有则增加set内容

	go func() {
		for d := range msgs_high {
			log.Printf("msgs_high.d: %s", d.Body)

			bytes := []byte(d.Body)
			broadcast(bytes)
		}
	}()
	go func() {
		for d := range msgs_normal {
			log.Printf("msgs_normal.new=%s", d.Body)
			bytes := []byte(d.Body)
			broadcast(bytes)
		}
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
	session_id := feedSessionStream.SessionID
	if map_topic_chanset[session_id] != nil {
		// broadcast 至特定的 session
		for val := range map_topic_chanset[session_id].Iterator().C {
			val.(chan []byte) <- bytes
		}
	}

}
