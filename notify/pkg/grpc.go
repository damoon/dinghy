//go:generate protoc -I .. --go_opt=paths=source_relative --go_out=plugins=grpc:pb service.proto

package notify

import (
	"sync"

	"gitlab.com/davedamoon/dinghy/notify/pkg/pb"
	"golang.org/x/net/context"
)

type GRPCServer struct {
	C *sync.Cond
}

func (s *GRPCServer) Listen(ctx context.Context, in *pb.Request) (*pb.Response, error) {
	s.C.L.Lock()
	s.C.Wait()
	s.C.L.Unlock()

	return &pb.Response{}, nil
}

func (s *GRPCServer) Notify(ctx context.Context, in *pb.Request) (*pb.Response, error) {
	s.C.L.Lock()
	s.C.Broadcast()
	s.C.L.Unlock()

	return &pb.Response{}, nil
}
