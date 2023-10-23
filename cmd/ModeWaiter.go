package cmd

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	pb "message/api/proto"
	"message/internal/config"
	"message/internal/service"
	"message/package/util_debug"
	"net"

	_ "net/http/pprof"
)

type server struct{}

func (s *server) UpdateSessionStatus(ctx context.Context, in *pb.UpdateSessionStatusReq) (*pb.UpdateSessionStatusResp, error) {
	log.Printf("ID: %s, Timestamp: %v , EvtType: %v ", in.SessionId, in.Timestamp, in.EvtType)

	service.OnNewStatus(in.SessionId, in.Timestamp, in.EvtType, in.Payload)
	return &pb.UpdateSessionStatusResp{ErrCode: 0, Msg: "OK UpdateSessionStatus"}, nil
}
func (s *server) FeedSessionStream(ctx context.Context, in *pb.FeedSessionStreamReq) (*pb.FeedSessionStreamResp, error) {
	log.Printf("SessionId: %s  ", in.SessionId)

	service.OnNewFeed(in.SessionId, in.Timestamp, in.Payload)
	return &pb.FeedSessionStreamResp{ErrCode: 0, Msg: "OK FeedSessionStream"}, nil
}

func MainModeWaiter() {
	debugOn, e := config.GetDebugMode()
	if e != nil {
		log.Fatal("GetDebugMode: ", e)
	}
	log.Println("debugOn : ", debugOn)

	if debugOn {
		addr, e := config.GetDebugPprofWaiter()
		if e != nil {
			log.Fatal("config e: ", e)
		}

		log.Println("GetDebugPprofWaiter addr: ", addr)
		go util_debug.InitPProf(addr)
	}

	v, _ := config.GetDependQueue()
	log.Println("GetDependQueue : ", v)
	service.InitProdQueue(v)

	v, _ = config.GetWaiterPortRpc()
	log.Println("GetWaiterPortRpc : ", v)
	listen, err := net.Listen("tcp", v)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
		return
	}

	ser := grpc.NewServer()
	pb.RegisterMessageServer(ser, &server{})

	reflection.Register(ser)
	if err := ser.Serve(listen); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
