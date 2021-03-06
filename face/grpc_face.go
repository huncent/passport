package face

import (
	"fmt"
	"net"

	"github.com/liuhengloveyou/passport/common"
	pb "github.com/liuhengloveyou/passport/face/pb"

	log "github.com/golang/glog"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type rpc struct{}

func (p *rpc) UserAdd(context.Context, *pb.UserInfo) (*pb.Reply, error) {
	return &pb.Reply{[]byte("this is demo message")}, nil

}

func (p *rpc) UserAuth(context.Context, *pb.Token) (*pb.Reply, error) {
	return nil, nil

}

func GrpcFace() {
	lis, err := net.Listen("tcp", common.ServConfig.Listen)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterPassportRpcServer(s, &rpc{})
	// Register reflection service on gRPC server.
	reflection.Register(s)

	fmt.Println("GO..." + common.ServConfig.Listen)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
