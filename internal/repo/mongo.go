package repo

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"message/api/domain"
	"message/internal/config"
	idomain "message/internal/domain"
	"strconv"
	"time"
)

const ( // QueryStatus
	ERR_OK           = 0
	ERR_NOT_FOUND    = 1404
	ERR_INVALID_PARA = 1500
)

var collection_log *mongo.Collection
var collection_status *mongo.Collection

var s_ctx context.Context

// DB static connection
func InitMongo(mongoDsn string) {
	var cancel context.CancelFunc
	s_ctx, cancel = context.WithTimeout(context.Background(), 2*time.Second)
	fmt.Println(cancel)
	//defer cancel()

	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(mongoDsn))
	if err != nil {
		log.Fatal(err)
	}

	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Ping OK")

	// biz code below . . .
	collection_log = client.Database(config.DATABASE).Collection(config.COLLECTION_LOG)
	collection_status = client.Database(config.DATABASE).Collection(config.COLLECTION_STATUS)
}

// pageID: 1,2,3,4....N
func GetSessionProfile(ctx context.Context, sessionID string, pageID int, pageSize int, detail bool) (error int, count int64, records []domain.BizMsg) {
	filter := bson.M{"session_id": sessionID}

	log.Printf("sessionID %s, pageID %d, size %d", sessionID, pageID, pageSize)

	cnt, _ := collection_log.CountDocuments(ctx, filter)
	fmt.Println("GetSessionProfile cnt: ", cnt)

	// 简略查询模式,直接返回记录数量
	if detail == false {
		return ERR_OK, cnt, nil
	}

	if cnt == 0 {
		return ERR_NOT_FOUND, 0, nil
	}

	findOpts := options.Find().SetSort(bson.D{{"timestamp", 1}}) // 1升序，-1降序
	if pageID > 0 && pageSize > 0 {
		findOpts.SetLimit(int64(pageSize))
		findOpts.SetSkip(int64(pageSize*pageID - pageSize))
	}
	//filter := bson.M{"session_id": session_id}
	findCursor, err := collection_log.Find(ctx, filter, findOpts)
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
		log.Printf("type: %d, timestamp:%d, payload:\n", value.Type, value.TimeStamp)
		//
		records = append(records, value)
	}

	cnt = int64(len(results))
	return ERR_OK, cnt, records
}

func GetSessionStatus(ctx context.Context, sessionID string, evtType string) (error int, payload string) {
	evtTypeI, _ := strconv.Atoi(evtType)
	filter := bson.M{"session_id": sessionID, "evt_type": evtTypeI}
	cnt, _ := collection_status.CountDocuments(ctx, filter)
	log.Println("GetSessionStatus cnt: ", cnt)

	log.Printf("session_id %s,  evt_type: %s ", sessionID, evtType)
	findCursor, err := collection_status.Find(ctx, filter)
	var results []bson.Raw
	if err = findCursor.All(ctx, &results); err != nil {
		log.Fatal(err)
	}

	for id, result := range results {
		log.Printf("id: %d, result: %v ", id, result)
		value := idomain.UpdateSessionStatus{}
		err2 := bson.Unmarshal([]byte(result), &value)
		log.Printf("Payload: %v ", value)

		if err2 != nil {
			panic(err)
		}
		return ERR_OK, value.Payload.(string)
	}

	return ERR_NOT_FOUND, ""
}
