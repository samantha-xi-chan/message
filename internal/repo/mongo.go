package repo

import (
	"context"
	"fmt"
	"github.com/message/api/domain"
	"github.com/message/config"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"time"
)

var collection *mongo.Collection

var s_ctx context.Context

const (
	ERR_OK           = 0
	ERR_INVALID_PARA = 1
)

// DB static connection
func InitMongo() {
	var cancel context.CancelFunc
	s_ctx, cancel = context.WithTimeout(context.Background(), 2*time.Second)
	fmt.Println(cancel)
	//defer cancel()
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		log.Fatal(err)
	}
	//
	//defer func() {
	//	if err = client.Disconnect(s_ctx); err != nil {
	//		panic(err)
	//	} else {
	//		fmt.Println("Connection to MongoDB closed.")
	//	}
	//}()

	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Ping OK")

	// biz code below . . .
	collection = client.Database(config.DATABASE).Collection(config.COLLECTION)
}

// pageID: 1,2,3,4....N
func GetSessionProfile(ctx context.Context, sessionID string, pageID int, pageSize int, detail bool) (error int, count int64, records []domain.BizMsg) {
	filter := bson.M{"session_id": sessionID}

	log.Printf("sessionID %s, pageID %d, size %d", sessionID, pageID, pageSize)

	cnt, _ := collection.CountDocuments(ctx, filter)
	fmt.Println("GetSessionProfile cnt: ", cnt)

	// 简略查询模式,直接返回记录数量
	if detail == false {
		return ERR_OK, cnt, nil
	}

	if cnt == 0 {
		return ERR_INVALID_PARA, 0, nil
	}

	findOpts := options.Find().SetSort(bson.D{{"timestamp", 1}}) // 1升序，-1降序
	if pageID > 0 && pageSize > 0 {
		findOpts.SetLimit(int64(pageSize))
		findOpts.SetSkip(int64(pageSize*pageID - pageSize))
	}
	//filter := bson.M{"session_id": session_id}
	findCursor, err := collection.Find(ctx, filter, findOpts)
	var results []bson.Raw
	if err = findCursor.All(ctx, &results); err != nil {
		log.Fatal(err)
	}

	for _, result := range results {
		value := domain.BizMsg{}
		err2 := bson.Unmarshal([]byte(result), &value)

		if err2 != nil {
			panic(err)
		}
		//fmt.Println("value: ", value)
		//fmt.Printf(" session_id: %s, timetamp: %d, payload: %v\n", value.SessionID, value.Timestamp, value.Payload)
		fmt.Printf("type: %d, timestamp:%d, payload:\n", value.Type, value.TimeStamp)
		//
		records = append(records, value)
	}

	cnt = int64(len(results))
	return ERR_OK, cnt, records
}
