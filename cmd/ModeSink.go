package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/message/config"
	"github.com/message/internal/domain"
	. "github.com/message/internal/service"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"time"
)

func MainModeSink() {
	//msgs_high := QueueConnInit(define.EXCHANGE_HIGH)
	msgs_normal, _ := QueueConnInit(config.EXCHANGE_NORMAL)

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		log.Fatal(err)
	}

	defer func() {
		if err = client.Disconnect(ctx); err != nil {
			panic(err)
		} else {
			fmt.Println("Connection to MongoDB closed.")
		}
	}()

	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Ping OK")

	// biz code below . . .
	collection := client.Database(config.DATABASE).Collection(config.COLLECTION)

	// 监听待推送的信息队列
	go func() {

		for d := range msgs_normal {
			log.Printf("接收消息=%s", d.Body)

			info := domain.FeedSessionStream{}
			json.Unmarshal(d.Body, &info)

			fmt.Printf("Unmarshal result: %v\n", info)

			// 落盘
			ctx, cancel = context.WithTimeout(context.Background(), 2*time.Second)
			defer cancel()
			res, _ := collection.InsertOne(ctx,
				bson.M{"session_id": info.SessionID, "timestamp": info.Timestamp, "payload": info.Payload})
			fmt.Printf("res.InsertedID: %v\n", res.InsertedID)

		}
	}()

	forever := make(chan bool)
	<-forever
}
