package dinghy

import (
	"context"
	"fmt"
	"io"
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
	ws, err := s.Upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		if _, ok := err.(websocket.HandshakeError); !ok {
			log.Println(err)
		}
		return
	}
	defer ws.Close()

	msg := make(chan []byte)
	defer close(msg)

	ctx := r.Context()

	go s.writer(ctx, ws, msg)
	reader(ws, msg)
}

func (s ServiceServer) CheckOrigin(r *http.Request) bool {
	if r.Header.Get("Origin") != s.FrontendURL {
		return false
	}

	return true
}

func reader(ws *websocket.Conn, msg chan<- []byte) {
	ws.SetReadLimit(512)
	ws.SetReadDeadline(time.Now().Add(pongWait))
	ws.SetPongHandler(func(string) error {
		ws.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	for {
		_, m, err := ws.ReadMessage()
		if err != nil && err != io.EOF {
			log.Printf("read from websocket: %v", err)
			break
		}

		msg <- m
	}
}

func (s ServiceServer) writer(ctx context.Context, ws *websocket.Conn, msg <-chan []byte) {
	pingTicker := time.NewTicker(pingPeriod)
	notify := s.listen(ctx)

	m := []byte{}
	var previous *Directory

	defer func() {
		pingTicker.Stop()
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

		case <-notify:
			if len(m) == 0 {
				continue
			}

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

	ws.SetWriteDeadline(time.Now().Add(writeWait))
	if err := ws.WriteJSON(listing); err != nil && err != io.EOF {
		return nil, fmt.Errorf("respond to websocket: %v", err)
	}

	return &listing, nil
}
