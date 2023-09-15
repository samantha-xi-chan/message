// client.go
package main

import (
	"context"
	"google.golang.org/grpc"
	"log"
	pb "message/api/proto"
	"message/internal/domain"
	"time"
)

const (
	address    = "localhost:10051"
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
		r, err := c.UpdateSessionStatus(ctx, &pb.UpdateSessionStatusReq{SessionId: SessionID6, Timestamp: time.Now().UnixNano() / 1e6, EvtType: domain.SESSION_START, Payload: "test_msg"})
		if err != nil {
			log.Fatalf("could not greet: %v", err)
		}
		log.Printf("ErrCode: %d, Msg: %s", r.ErrCode, r.Msg)
	}

	// 任务过程中的细节变化
	for i := 0; i <= 2; i++ {
		tick := time.Now().UnixNano() / 1e6
		val := "test_msg_2"
		r, err := c.FeedSessionStream(ctx, &pb.FeedSessionStreamReq{SessionId: SessionID6, Timestamp: tick, Payload: val})
		if err != nil {
			log.Fatalf("could not FeedSessionStream: %v", err)
		}
		log.Printf("ErrCode: %d, Msg: %s, posted:%s ", r.ErrCode, r.Msg, val)
	}

	// 任务结束通知
	{
		val := "exit 2"
		r, err := c.UpdateSessionStatus(ctx, &pb.UpdateSessionStatusReq{SessionId: SessionID6, Timestamp: time.Now().UnixNano() / 1e6, EvtType: domain.SESSION_END, Payload: val})
		if err != nil {
			log.Fatalf("could not greet: %v", err)
		}
		log.Printf("ErrCode: %d, Msg: %s", r.ErrCode, r.Msg)
	}

}
