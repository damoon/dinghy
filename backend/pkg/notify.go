package dinghy

import (
	"context"
	"io"
	"log"
	"os"
	"time"

	"gitlab.com/davedamoon/dinghy/notify/pkg/pb"
	"google.golang.org/grpc"
)

func notify() {
	// Set up a connection to the server.
	conn, err := grpc.Dial("notify:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewNotifierClient(conn)

	// Contact the server and print out its response.
	name := "world"
	if len(os.Args) > 1 {
		name = os.Args[1]
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	r, err := c.Listen(ctx, &pb.Request{Name: name})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
}

func listen() (<-chan struct{}, io.Closer) {
	// Set up a connection to the server.
	conn, err := grpc.Dial("notify:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewNotifierClient(conn)

	// Contact the server and print out its response.
	name := "world"
	if len(os.Args) > 1 {
		name = os.Args[1]
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	r, err := c.Listen(ctx, &pb.Request{Name: name})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
}
