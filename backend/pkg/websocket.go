package dinghy

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"reflect"
	"time"

	//	"github.com/go-redis/redis"
	"github.com/gorilla/websocket"
)

const (
	writeWait     = 10 * time.Second
	pongWait      = 60 * time.Second
	pingPeriod    = (pongWait * 9) / 10
	refreshPeriod = 2 * time.Second
)

func (s ServiceServer) serveWs(w http.ResponseWriter, r *http.Request) {
	ws, err := s.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		if _, ok := err.(websocket.HandshakeError); !ok {
			log.Println(err)
		}
		return
	}

	msg := make(chan []byte)
	defer close(msg)

	go s.writer(ws, msg)
	reader(ws, msg)
}

func reader(ws *websocket.Conn, msg chan<- []byte) {
	defer ws.Close()

	ws.SetReadLimit(512)
	ws.SetReadDeadline(time.Now().Add(pongWait))
	ws.SetPongHandler(func(string) error {
		ws.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	for {
		// TODO: check and ignore closed
		_, m, err := ws.ReadMessage()
		if err != nil {
			log.Printf("read from websocket: %v", err)
			break
		}

		msg <- m
	}
}

func (s ServiceServer) writer(ws *websocket.Conn, msg <-chan []byte) {
	pingTicker := time.NewTicker(pingPeriod)
	fileTicker := time.NewTicker(refreshPeriod)
	m := []byte{}
	var previous *Directory

	defer func() {
		pingTicker.Stop()
		fileTicker.Stop()
		ws.Close()
	}()

	for {
		select {
		case m = <-msg:
			cur, err := s.sendUpdate(ws, nil, string(m))
			if err != nil {
				log.Println(err)
				return
			}

			previous = cur

		//		case <-fileTicker.C:
		case <-notify:
			if len(m) == 0 {
				continue
			}

			// TODO failed notify connection should trigger reload

			//			client := redis.NewClient(&redis.Options{
			//				Addr:     "redis:6379",
			//				Password: "redis123",
			//			})
			//			defer client.Close()
			log.Println("notify")
			cur, err := s.sendUpdate(ws, previous, string(m))
			if err != nil {
				log.Println(err)
				return
			}

			previous = cur

		case <-pingTicker.C:
			ws.SetWriteDeadline(time.Now().Add(writeWait))
			if err := ws.WriteMessage(websocket.PingMessage, []byte{}); err != nil {
				log.Println(err)
				return
			}
		}
	}
}

func (s ServiceServer) sendUpdate(ws *websocket.Conn, previous *Directory, path string) (*Directory, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	listing, err := s.Storage.list(ctx, path)
	if err != nil {
		return nil, fmt.Errorf("list %s: %v", path, err)
	}

	if reflect.DeepEqual(previous, &listing) {
		return &listing, nil
	}

	// TODO: check and ignore closed
	ws.SetWriteDeadline(time.Now().Add(writeWait))
	if err := ws.WriteJSON(listing); err != nil {
		return nil, fmt.Errorf("respond to websocket: %v", err)
	}

	return &listing, nil
}
