// client.go
package main

import (
	"context"
	"encoding/json"
	"google.golang.org/grpc"
	"log"
	apidomain "message/api/domain"
	pb "message/api/proto"
	"message/internal/domain"
	"time"
)

const (
	address = "localhost:10051"

	SessionID6 = "test"
)

func main() {
	for i := 0; i < 100; i++ {
		send()
		time.Sleep(time.Millisecond * 10)
	}
}

func send() {
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewMessageClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*200000)
	defer cancel()

	// 任务启动通知
	{
		address := apidomain.SessionProperty{
			Src: "22:22:22",
			Dst: "33:33:33",
		}
		jBiz, err := json.Marshal(address)

		r, err := c.UpdateSessionStatus(ctx, &pb.UpdateSessionStatusReq{SessionId: SessionID6, Timestamp: time.Now().UnixNano() / 1e6, EvtType: domain.SESSION_START, Payload: string(jBiz)})
		if err != nil {
			log.Fatalf("could not greet: %v", err)
		}
		log.Printf("ErrCode: %d, Msg: %s", r.ErrCode, r.Msg)
	}

	// 任务过程中的细节变化
	for i := 0; i <= 2; i++ {
		tick := time.Now().UnixNano() / 1e6
		biz := apidomain.BizMsg{
			Type:      apidomain.CASE_CHANGE, // 1001表示用例变化 ， 1002表示用例内的报文日志
			TimeStamp: tick,
			Payload: apidomain.CaseChange{
				CaseID:         "Case1001",
				CaseChangeType: 0,
				CaseAttr:       "att",
			},
		}
		jBiz, err := json.Marshal(biz)

		r, err := c.FeedSessionStream(ctx, &pb.FeedSessionStreamReq{SessionId: SessionID6, Timestamp: tick, Payload: string(jBiz)})
		if err != nil {
			log.Fatalf("could not FeedSessionStream: %v", err)
		}
		log.Printf("ErrCode: %d, Msg: %s, posted:%s ", r.ErrCode, r.Msg, string(jBiz))
		//time.Sleep(time.Millisecond * 20)
	}

	// 任务结束通知
	{
		r, err := c.UpdateSessionStatus(ctx, &pb.UpdateSessionStatusReq{SessionId: SessionID6, Timestamp: time.Now().UnixNano() / 1e6, EvtType: domain.SESSION_END, Payload: "{\"exit_code\":0}"})
		if err != nil {
			log.Fatalf("could not greet: %v", err)
		}
		log.Printf("ErrCode: %d, Msg: %s", r.ErrCode, r.Msg)
	}

}
