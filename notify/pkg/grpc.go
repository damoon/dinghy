//go:generate protoc -I ../. --go_out=plugins=grpc:../pkg/pb service.proto

package notify

import (
	"gitlab.com/davedamoon/dinghy/notify/pkg/pb"
	"golang.org/x/net/context"
)

type GRPCServer struct{}

func (s *GRPCServer) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
	return &pb.HelloReply{Message: "Hello " + in.Name}, nil
}
