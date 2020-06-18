package dinghy

import (
	"context"
	"log"
	"time"

	"gitlab.com/davedamoon/dinghy/backend/pkg/pb"
)

func (s *ServiceServer) notify() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	_, err := s.NotifierClient.Notify(ctx, &pb.Request{})
	if err != nil {
		log.Printf("could not notify: %v", err)
	}
}

func (s *ServiceServer) listen(ctx context.Context) <-chan struct{} {
	ch := make(chan struct{})

	go func() {
		for {
			_, err := s.NotifierClient.Listen(ctx, &pb.Request{})

			select {
			case <-ctx.Done():
				return
			default:
			}

			if err != nil {
				log.Printf("could not listen: %v", err)
				time.Sleep(time.Second)
			}

			ch <- struct{}{}
		}
	}()

	return ch
}
