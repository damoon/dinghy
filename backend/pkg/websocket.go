package dinghy

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

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

	defer func() {
		pingTicker.Stop()
		fileTicker.Stop()
		ws.Close()
	}()

	for {
		select {
		case m = <-msg:
			err := s.sendUpdate(ws, string(m))
			if err != nil {
				log.Println(err)
				return
			}

		case <-fileTicker.C:
			if len(m) != 0 {
				err := s.sendUpdate(ws, string(m))
				if err != nil {
					log.Println(err)
					return
				}
			}

		case <-pingTicker.C:
			ws.SetWriteDeadline(time.Now().Add(writeWait))
			if err := ws.WriteMessage(websocket.PingMessage, []byte{}); err != nil {
				log.Println(err)
				return
			}
		}
	}
}

func (s ServiceServer) sendUpdate(ws *websocket.Conn, path string) error {
	ws.SetWriteDeadline(time.Now().Add(writeWait))

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	l, err := s.Storage.list(ctx, path)
	if err != nil {
		return fmt.Errorf("list %s: %v", path, err)
	}

	// TODO: check and ignore closed
	if err := ws.WriteJSON(l); err != nil {
		return fmt.Errorf("respond to websocket: %v", err)
	}

	return nil
}
