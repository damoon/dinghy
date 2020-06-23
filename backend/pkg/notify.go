package dinghy

import (
	"context"
	"log"
	"time"

	"gitlab.com/davedamoon/dinghy/backend/pkg/pb"
)

type NotifyAdapter struct {
	NotifierClient pb.NotifierClient
}

func (n *NotifyAdapter) notify(ctx context.Context) {
	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()
	_, err := n.NotifierClient.Notify(ctx, &pb.Request{})
	if err != nil {
		log.Printf("could not notify: %v", err)
	}
}

func (n *NotifyAdapter) listen(ctx context.Context) <-chan struct{} {
	ch := make(chan struct{})

	go func() {
		for {
			_, err := n.NotifierClient.Listen(ctx, &pb.Request{})

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
