package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/Clouditera/message/config"
	"github.com/Clouditera/message/internal"
	"github.com/Clouditera/message/internal/domain"
	. "github.com/Clouditera/message/internal/service"
	"net/http"

	//"github.com/sirupsen/fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"time"

	_ "net/http/pprof"
)

type LogEntry struct {
	SessionID string
	Timestamp int64
	Payload   string
}

const TIMEOUT_SECOND = 99999
const GC_INTERVAL_SECOND = 600

func MainModeSink() {
	go func() {
		log.Println(http.ListenAndServe("0.0.0.0:6080", nil))
	}()

	msgs_high, _ := QueueConnInit(config.EXCHANGE_HIGH)
	msgs_normal, _ := QueueConnInit(config.EXCHANGE_NORMAL)

	ctx, cancel := context.WithTimeout(context.Background(), TIMEOUT_SECOND*time.Second)
	defer cancel()
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(internal.MONGO_URL).SetMaxPoolSize(100))
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
	collection_log := client.Database(config.DATABASE).Collection(config.COLLECTION_LOG)
	res, err := collection_log.Indexes().CreateOne(context.Background(), mongo.IndexModel{
		Keys: bson.D{
			{Key: "session_id", Value: 1},
		},
	})
	if err != nil {
		log.Fatal("Indexes().CreateOne err: ", err)
	}
	log.Println(res)
	res, err = collection_log.Indexes().CreateOne(context.Background(), mongo.IndexModel{
		Keys: bson.D{
			{Key: "timestamp", Value: 1},
		},
	})
	if err != nil {
		log.Fatal("Indexes().CreateOne err: ", err)
	}
	log.Println(res)
	res, err = collection_log.Indexes().CreateOne(context.Background(), mongo.IndexModel{
		Keys: bson.D{
			{Key: "session_id", Value: 1},
			{Key: "timestamp", Value: 1},
		},
	})
	if err != nil {
		log.Fatal("Indexes().CreateOne err: ", err)
	}
	log.Println(res)

	/*
		go func() {
			for d := range msgs_normal {
				log.Printf("recv msgs_normal = %s", d.Body)

				info := domain.FeedSessionStream{}
				json.Unmarshal(d.Body, &info)

				//fmt.Printf("Unmarshal result: %v\n", info)

				// 落盘
				ctx, cancel = context.WithTimeout(context.Background(), TIMEOUT_SECOND*time.Second)
				defer cancel()
				res, _ := collection_log.InsertOne(ctx,
					bson.M{"session_id": info.SessionID, "timestamp": info.Timestamp, "payload": info.Payload, "deleted": false})
				fmt.Printf("res.InsertedID: %v\n", res.InsertedID)
			}
		}()
	*/

	go func() {
		batchData := []interface{}{}

		const batchSize = 100
		const timeout = 5 * time.Second
		timer := time.NewTimer(timeout)
		for {
			select {
			case d := <-msgs_normal:
				info := domain.FeedSessionStream{}
				json.Unmarshal(d.Body, &info)
				batchData = append(batchData, bson.M{"session_id": info.SessionID, "timestamp": info.Timestamp, "payload": info.Payload, "deleted": false})
				if len(batchData) >= batchSize {
					collection_log.InsertMany(ctx, batchData)
					batchData = nil
					timer.Reset(timeout)
				}
			case <-timer.C:
				if len(batchData) > 0 {
					collection_log.InsertMany(ctx, batchData)
					batchData = nil
				}
				timer.Reset(timeout)
			}
		}
	}()

	collection_status := client.Database(config.DATABASE).Collection(config.COLLECTION_STATUS)
	res, err = collection_status.Indexes().CreateOne(context.Background(), mongo.IndexModel{
		Keys: bson.D{
			{Key: "session_id", Value: 1},
		},
	})
	if err != nil {
		log.Fatal("Indexes().CreateOne err: ", err)
	}
	log.Println(res)
	go func() {
		/*
			for d := range msgs_high {
				log.Printf("recv msgs_high =%s", d.Body)

				info := domain.UpdateSessionStatus{}
				json.Unmarshal(d.Body, &info)

				//fmt.Printf("Unmarshal result: %v\n", info)

				ctx, cancel = context.WithTimeout(context.Background(), TIMEOUT_SECOND*time.Second)
				defer cancel()
				res, _ := collection_status.InsertOne(ctx,
					bson.M{"session_id": info.SessionID, "timestamp": info.Timestamp, "evt_type": info.EvtType, "payload": info.Payload, "deleted": false})
				fmt.Printf("res.InsertedID: %v\n", res.InsertedID)
			}
		*/
		batchData := []interface{}{}

		const batchSize = 100
		const timeout = 1 * time.Second
		timer := time.NewTimer(timeout)
		for {
			select {
			case d := <-msgs_high:
				info := domain.UpdateSessionStatus{}
				json.Unmarshal(d.Body, &info)
				batchData = append(batchData, bson.M{"session_id": info.SessionID, "timestamp": info.Timestamp, "evt_type": info.EvtType, "payload": info.Payload, "deleted": false})
				if len(batchData) >= batchSize {
					collection_status.InsertMany(ctx, batchData)
					batchData = nil
					timer.Reset(timeout)
				}
			case <-timer.C:
				if len(batchData) > 0 {
					collection_status.InsertMany(ctx, batchData)
					batchData = nil
				}
				timer.Reset(timeout)
			}
		}

	}()

	fmt.Println("starting cron job")
	//c := cron.New()
	//e := c.AddFunc("*/3 * * * *", func() {
	collection_log2 := client.Database(config.DATABASE).Collection(config.COLLECTION_LOG)
	for false { // disable it
		fmt.Print("cron func() running ")
		var MAX_BUF_SIZE = 5000

		ctx, cancel = context.WithTimeout(context.Background(), TIMEOUT_SECOND*time.Second)
		defer cancel()

		// 构建聚合管道，按 session_id 分组，统计记录数量，并筛选出记录数量大于 5000 的 session_id
		pipeline := []bson.M{
			{
				"$group": bson.M{
					"_id":   "$session_id",
					"count": bson.M{"$sum": 1},
				},
			},
			{
				"$match": bson.M{
					"count": bson.M{"$gt": MAX_BUF_SIZE},
				},
			},
		}

		// 执行聚合操作
		cursor, err := collection_log2.Aggregate(ctx, pipeline)
		if err != nil {
			fmt.Print("Aggregate err: ", err)
		}

		// 遍历结果
		var sessionIDs []string
		for cursor.Next(ctx) {
			var result bson.M
			err := cursor.Decode(&result)
			if err != nil {
				fmt.Println("Decode err: ", err)
			}

			sessionID := result["_id"].(string)
			fmt.Println("trying to cleanup sessionID: ", sessionID)

			//
			{
				findFilter := bson.M{"session_id": sessionID}
				findOptions := options.Find().SetSort(bson.D{{"timestamp", -1}})
				findCursor, err := collection_log2.Find(ctx, findFilter, findOptions)
				if err != nil {
					fmt.Println("Find err: ", err)
				}

				var recordsToDelete []LogEntry
				for findCursor.Next(ctx) {
					var logEntry LogEntry
					err := findCursor.Decode(&logEntry)
					if err != nil {
						fmt.Println("Decode err: ", err)
					}

					recordsToDelete = append(recordsToDelete, logEntry)
				}

				// 只保留前 5000 条记录
				if len(recordsToDelete) > MAX_BUF_SIZE {
					recordsToDelete = recordsToDelete[MAX_BUF_SIZE:]
				} else {
					recordsToDelete = nil
				}

				// 删除 session_id 内的记录
				for _, logEntry := range recordsToDelete {
					deleteFilter := bson.M{"_id": logEntry.SessionID}
					_, err := collection_log2.DeleteOne(ctx, deleteFilter)
					if err != nil {
						fmt.Println("DeleteOne err: ", err)
					}
				}

				fmt.Printf("Deleted %d records for session ID %s\n", len(recordsToDelete), sessionID)
			}
		}

		fmt.Println("Session IDs with more than MAX_BUF_SIZE records:", sessionIDs)
		time.Sleep(time.Second * GC_INTERVAL_SECOND)
	} // end of for

	select {}
}
