package dinghy

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"reflect"
	"time"

	"github.com/gorilla/websocket"
)

const (
	writeWait  = 10 * time.Second
	pongWait   = 60 * time.Second
	pingPeriod = (pongWait * 9) / 10
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

func reader(ws *websocket.Conn, msg chan<- []byte) {
	ws.SetReadLimit(512)

	err := ws.SetReadDeadline(time.Now().Add(pongWait))
	if err != nil {
		log.Printf("websocket reader: %v", err)
		return
	}

	ws.SetPongHandler(func(string) error {
		return ws.SetReadDeadline(time.Now().Add(pongWait))
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
	notify := s.Notify.listen(ctx)

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
			err := ws.SetWriteDeadline(time.Now().Add(writeWait))
			if err != nil {
				log.Println(err)
				return
			}

			err = ws.WriteMessage(websocket.PingMessage, []byte{})
			if err != nil {
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

	err = ws.SetWriteDeadline(time.Now().Add(writeWait))
	if err != nil {
		return nil, fmt.Errorf("set write deadline: %v", err)
	}

	err = ws.WriteJSON(listing)
	if err != nil && err != io.EOF {
		return nil, fmt.Errorf("respond to websocket: %v", err)
	}

	return &listing, nil
}
