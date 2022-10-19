package cmd

import (
	"encoding/json"
	"fmt"
	mapset "github.com/deckarep/golang-set"
	"github.com/gorilla/websocket"
	api "github.com/message/api/domain"
	"github.com/message/config"
	"github.com/message/internal/domain"
	. "github.com/message/internal/service"
	"log"
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
	w_chan := make(chan []byte, 99999)

	go func() {
		for {
			msg, ok := <-w_chan
			if ok {
				err = conn.WriteMessage(websocket.TextMessage, msg)
				if err != nil {
					log.Println("Error during message writing:", err)
					break
				}
			}
		}
	}()

	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			log.Println("Error during message reading:", err)
			break
		}
		log.Printf("Received: %s", message)

		// 解析是否是订阅,如果是订阅则新增映射
		req := api.WsReq{}
		json.Unmarshal(message, &req)

		if req.Type == config.TYPE_SUBSCRIBE && req.Version == config.WS_PROTO_VER {
			//sub := api.Subscribe{}

			// unmarshal it (usually after receiving bytes from somewhere)
			sub := &api.Subscribe{}
			wsReq := api.WsReq{Payload: sub}
			json.Unmarshal(message, &wsReq)

			//bytes := []byte(fmt.Sprintf("%v", req.Payload.(interface{})))
			//log.Printf("bytes: %s", string(bytes))
			//json.Unmarshal(bytes, &sub)
			topic = sub.Topic
			log.Printf("topic: %s", topic)
			if map_topic_chanset[topic] == nil {
				// 新建一条 map的KV记录
				log.Printf("new map KV: %s", topic)
				map_topic_chanset[topic] = mapset.NewSet()
			}
			map_topic_chanset[topic].Add(w_chan)

			// 暂时不需要其他处理 其实可以直接丢弃
			//r_chan <- message
		}
	}

	// 移除当前 topic 对本次发送缓冲区的映射
	map_topic_chanset[topic].Remove(w_chan)

	fmt.Println("socketHandlerB ending...")
	//forever := make(chan bool)
	//<-forever
}

func MainModeNotify() {
	msgs_high, _ := QueueConnInit(config.EXCHANGE_HIGH)
	msgs_normal, _ := QueueConnInit(config.EXCHANGE_NORMAL)

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
	log.Fatal(http.ListenAndServe(":18080", nil))
}

func broadcast(bytes []byte) {
	var feedSessionStream domain.FeedSessionStream
	err := json.Unmarshal(bytes, &feedSessionStream)
	if err != nil {
		fmt.Println("Error", err)
	} else {
		fmt.Println("feedSessionStream: ", feedSessionStream)
	}
	session_id := feedSessionStream.SessionID
	if map_topic_chanset[session_id] != nil {
		// broadcast 至特定的 session
		for val := range map_topic_chanset[session_id].Iterator().C {
			val.(chan []byte) <- bytes
		}
	}

}
