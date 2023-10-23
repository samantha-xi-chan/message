package cmd

import (
	"context"
	"encoding/json"
	"message/internal/repo"
	"message/internal/service"
	"message/package/util_debug"

	"log"
	"message/internal/config"
	"message/internal/domain"
	//. "message/internal/service"
	"time"

	_ "net/http/pprof"
)

type LogEntry struct {
	SessionID string
	Timestamp int64
	Payload   string
}

//const TIMEOUT_SECOND = 99999
//const GC_INTERVAL_SECOND = 600

func MainModeSink() {
	debugOn, e := config.GetDebugMode()
	if e != nil {
		log.Fatal("GetDebugMode: ", e)
	}

	if debugOn {
		addr, e := config.GetDebugPprofSink()
		if e != nil {
			log.Fatal("config e: ", e)
		}

		log.Println("GetDebugPprofSink addr: ", addr)
		go util_debug.InitPProf(addr)
	}

	redisDsn, e := config.GetDependRedisDsn()
	if e != nil {
		log.Fatal("GetDependRedisDsn: ", e)
	}
	log.Println("redisDsn: ", redisDsn)

	//storeMaxCount, e := config.GetStoreMaxCount()
	//if e != nil {
	//	log.Fatal("GetStoreMaxCount: ", e)
	//}
	storeMaxCount := int64(100)
	log.Println("storeMaxCount: ", storeMaxCount)

	e = repo.InitRedis(context.Background(), redisDsn, storeMaxCount, 0)
	if e != nil {
		log.Fatal("InitRedis: ", e)
	}

	v, _ := config.GetDependQueue()
	msgs_high, _ := service.QueueConnInit(v, config.EXCHANGE_HIGH)
	msgs_normal, _ := service.QueueConnInit(v, config.EXCHANGE_NORMAL)

	//ctx, cancel := context.WithTimeout(context.Background(), TIMEOUT_SECOND*time.Second)
	//defer cancel()

	//val, _ := config.GetDependMongo()
	//client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(val).SetMaxPoolSize(100))
	//if err != nil {
	//	log.Fatal(err)
	//}

	//defer func() {
	//	if err = client.Disconnect(ctx); err != nil {
	//		panic(err)
	//	} else {
	//		fmt.Println("Connection to MongoDB closed.")
	//	}
	//}()
	//
	//err = client.Ping(context.TODO(), nil)
	//if err != nil {
	//	log.Fatal(err)
	//}
	//fmt.Println("Ping OK")

	// biz code below . . .
	//collection_log := client.Database(config.DATABASE).Collection(config.COLLECTION_LOG)
	//res, err := collection_log.Indexes().CreateOne(context.Background(), mongo.IndexModel{
	//	Keys: bson.D{
	//		{Key: "session_id", Value: 1},
	//	},
	//})
	//if err != nil {
	//	log.Fatal("Indexes().CreateOne err: ", err)
	//}
	//log.Println(res)
	//res, err = collection_log.Indexes().CreateOne(context.Background(), mongo.IndexModel{
	//	Keys: bson.D{
	//		{Key: "timestamp", Value: 1},
	//	},
	//})
	//if err != nil {
	//	log.Fatal("Indexes().CreateOne err: ", err)
	//}
	//log.Println(res)
	//res, err = collection_log.Indexes().CreateOne(context.Background(), mongo.IndexModel{
	//	Keys: bson.D{
	//		{Key: "session_id", Value: 1},
	//		{Key: "timestamp", Value: 1},
	//	},
	//})
	//if err != nil {
	//	log.Fatal("Indexes().CreateOne err: ", err)
	//}
	//log.Println(res)

	msgs_normal_buffered := make(chan []byte, 10)
	go func() {
		for dd := range msgs_normal {
			msgs_normal_buffered <- dd.Body
		}
		close(msgs_normal_buffered)
	}()

	go func() {
		//batchData := []interface{}{}

		timer := time.NewTimer(config.FLUSH_BUF_TIMEOUT_SEC * time.Second)
		for i := 0; ; i++ {
			select {
			case d, ok := <-msgs_normal_buffered:
				if !ok {
					log.Println("msgs_high_buffered closed, exiting loop")
					return
				}

				info := domain.FeedSessionStream{}
				json.Unmarshal(d, &info)

				e := repo.GetRedisMgr().NewLog(context.Background(), i%100 == 0, info.SessionID, string(d))
				if e != nil {
					log.Println("GetRedisMgr().NewLog: ", e)
				}

				//batchData = append(batchData, bson.M{"session_id": info.SessionID, "timestamp": info.Timestamp, "payload": info.Payload, "deleted2": false})
				//if len(batchData) >= config.BATCH_SIZE {
				//	collection_log.InsertMany(ctx, batchData)
				//	log.Println("batchData.size batchSize log: ", len(batchData))
				//	batchData = nil
				//	timer.Reset(config.FLUSH_BUF_TIMEOUT_SEC * time.Second)
				//}
			case <-timer.C:
				//if len(batchData) > 0 {
				//	collection_log.InsertMany(ctx, batchData)
				//	log.Println("batchData.size timer log: ", len(batchData))
				//	batchData = nil
				//}
				timer.Reset(config.FLUSH_BUF_TIMEOUT_SEC * time.Second)
			}
		}
	}()

	//collection_status := client.Database(config.DATABASE).Collection(config.COLLECTION_STATUS)
	//res, err = collection_status.Indexes().CreateOne(context.Background(), mongo.IndexModel{
	//	Keys: bson.D{
	//		{Key: "session_id", Value: 1},
	//	},
	//})
	//if err != nil {
	//	log.Fatal("Indexes().CreateOne err: ", err)
	//}
	//log.Println(res)

	msgs_high_buffered := make(chan []byte, 10)
	go func() {
		for dd := range msgs_high {
			msgs_high_buffered <- dd.Body
		}
		close(msgs_high_buffered)
	}()
	go func() {
		//batchData := []interface{}{}
		timer := time.NewTimer(config.FLUSH_BUF_TIMEOUT_SEC * time.Second)
		for {
			select {
			case _, ok := <-msgs_high_buffered:
				if !ok {
					log.Println("msgs_high_buffered closed, exiting loop")
					return
				}

				//info := domain.UpdateSessionStatus{}
				//json.Unmarshal(d, &info)
				//batchData = append(batchData, bson.M{"session_id": info.SessionID, "timestamp": info.Timestamp, "evt_type": info.EvtType, "payload": info.Payload, "deleted2": false})
				//if len(batchData) >= config.BATCH_SIZE {
				//	collection_status.InsertMany(ctx, batchData)
				//	log.Println("batchData.size batchSize status: ", len(batchData))
				//	batchData = nil
				//	timer.Reset(config.FLUSH_BUF_TIMEOUT_SEC * time.Second)
				//}
			case <-timer.C:
				//if len(batchData) > 0 {
				//	collection_status.InsertMany(ctx, batchData)
				//	log.Println("batchData.size timer status: ", len(batchData))
				//	batchData = nil
				//}
				timer.Reset(config.FLUSH_BUF_TIMEOUT_SEC * time.Second)
			}
		}

	}()

	//fmt.Println("starting cron job")
	//c := cron.New()
	//e := c.AddFunc("*/3 * * * *", func() {
	//collection_log2 := client.Database(config.DATABASE).Collection(config.COLLECTION_LOG)
	//for false { // disable it
	//	fmt.Print("cron func() running ")
	//	var MAX_BUF_SIZE = 5000
	//
	//	ctx, cancel = context.WithTimeout(context.Background(), TIMEOUT_SECOND*time.Second)
	//	defer cancel()
	//
	//	// 构建聚合管道，按 session_id 分组，统计记录数量，并筛选出记录数量大于 5000 的 session_id
	//	pipeline := []bson.M{
	//		{
	//			"$group": bson.M{
	//				"_id":   "$session_id",
	//				"count": bson.M{"$sum": 1},
	//			},
	//		},
	//		{
	//			"$match": bson.M{
	//				"count": bson.M{"$gt": MAX_BUF_SIZE},
	//			},
	//		},
	//	}
	//
	//	// 执行聚合操作
	//	cursor, err := collection_log2.Aggregate(ctx, pipeline)
	//	if err != nil {
	//		fmt.Print("Aggregate err: ", err)
	//	}
	//
	//	// 遍历结果
	//	var sessionIDs []string
	//	for cursor.Next(ctx) {
	//		var result bson.M
	//		err := cursor.Decode(&result)
	//		if err != nil {
	//			fmt.Println("Decode err: ", err)
	//		}
	//
	//		sessionID := result["_id"].(string)
	//		fmt.Println("trying to cleanup sessionID: ", sessionID)
	//
	//		//
	//		{
	//			findFilter := bson.M{"session_id": sessionID}
	//			findOptions := options.Find().SetSort(bson.D{{"timestamp", -1}})
	//			findCursor, err := collection_log2.Find(ctx, findFilter, findOptions)
	//			if err != nil {
	//				fmt.Println("Find err: ", err)
	//			}
	//
	//			var recordsToDelete []LogEntry
	//			for findCursor.Next(ctx) {
	//				var logEntry LogEntry
	//				err := findCursor.Decode(&logEntry)
	//				if err != nil {
	//					fmt.Println("Decode err: ", err)
	//				}
	//
	//				recordsToDelete = append(recordsToDelete, logEntry)
	//			}
	//
	//			// 只保留前 5000 条记录
	//			if len(recordsToDelete) > MAX_BUF_SIZE {
	//				recordsToDelete = recordsToDelete[MAX_BUF_SIZE:]
	//			} else {
	//				recordsToDelete = nil
	//			}
	//
	//			// 删除 session_id 内的记录
	//			for _, logEntry := range recordsToDelete {
	//				deleteFilter := bson.M{"_id": logEntry.SessionID}
	//				_, err := collection_log2.DeleteOne(ctx, deleteFilter)
	//				if err != nil {
	//					fmt.Println("DeleteOne err: ", err)
	//				}
	//			}
	//
	//			fmt.Printf("Deleted %d records for session ID %s\n", len(recordsToDelete), sessionID)
	//		}
	//	}
	//
	//	fmt.Println("Session IDs with more than MAX_BUF_SIZE records:", sessionIDs)
	//	time.Sleep(time.Second * GC_INTERVAL_SECOND)
	//} // end of for

	select {}
}
