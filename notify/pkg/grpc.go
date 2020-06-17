//go:generate protoc -I ../. --go_out=plugins=grpc:../pkg/pb service.proto

package notify

import (
	"log"
	"sync"

	"gitlab.com/davedamoon/dinghy/notify/pkg/pb"
	"golang.org/x/net/context"
)

type GRPCServer struct {
	C *sync.Cond
}

func (s *GRPCServer) Listen(ctx context.Context, in *pb.Request) (*pb.Response, error) {
	log.Println("Listen")
	s.C.L.Lock()
	s.C.Wait()
	s.C.L.Unlock()
	log.Println("Listen Done")

	return &pb.Response{Message: "Hello " + in.Name}, nil
}

func (s *GRPCServer) Notify(ctx context.Context, in *pb.Request) (*pb.Response, error) {
	log.Println("Notify")
	s.C.L.Lock()
	s.C.Broadcast()
	s.C.L.Unlock()

	return &pb.Response{Message: "Hello " + in.Name}, nil
}
