package dinghy

import (
	"context"
	"fmt"
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

func (s ServiceServer) serveWs(w http.ResponseWriter, r *http.Request) error {
	ws, err := s.Upgrader.Upgrade(w, r, nil)
	if err != nil {
		if _, ok := err.(websocket.HandshakeError); !ok {
			return err
		}
	}
	defer ws.Close()

	msg := make(chan string)
	defer close(msg)

	ctx := r.Context()

	go s.writer(ctx, ws, msg)
	s.reader(ctx, ws, msg)

	return nil
}

func (s ServiceServer) reader(ctx context.Context, ws *websocket.Conn, msg chan<- string) {
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
		if err != nil {
			e, ok := err.(*websocket.CloseError)
			if !ok {
				log.Printf("read from websocket: %v", err)
				break
			}
			if e.Code == websocket.CloseGoingAway {
				break
			}
			log.Printf("read from websocketdssss: %v", err)
			break
		}

		switch string(m[0:3]) {
		case "cd ":
			msg <- string(m[3:])
		case "ex ":
			go func(path string) {
				err := s.unzip(ctx, path)
				if err != nil {
					log.Printf("extract %s: %v", path, err)
				}
				s.Notify.notify(ctx)
			}(string(m[3:]))
		case "rm ":
			path := string(m[3:])

			err := s.Storage.deleteRecursive(ctx, path)
			if err != nil {
				log.Printf("deleting %s: %v", path, err)
			}

			s.Notify.notify(ctx)
		}

	}
}

func (s ServiceServer) writer(ctx context.Context, ws *websocket.Conn, msg <-chan string) {
	pingTicker := time.NewTicker(pingPeriod)
	notify := s.Notify.listen(ctx)

	path := ""
	var previous *Directory

	defer func() {
		pingTicker.Stop()
	}()

	for {
		select {

		case <-ctx.Done():
			return

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
		case <-notify:
			if len(path) == 0 {
				continue
			}

			cur, err := s.sendUpdate(ctx, ws, previous, path)
			if err != nil {
				log.Println(err)
				return
			}

			previous = cur
		case path = <-msg:
			if len(path) == 0 {
				return
			}

			cur, err := s.sendUpdate(ctx, ws, nil, path)
			if err != nil {
				log.Println(err)
				return
			}

			previous = cur
		}
	}
}

func (s ServiceServer) sendUpdate(ctx context.Context, ws *websocket.Conn, previous *Directory, path string) (*Directory, error) {
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
	if err != nil {
		if err == websocket.ErrCloseSent {
			return nil, nil
		}
		return nil, fmt.Errorf("respond to websocket: %v", err)
	}

	return &listing, nil
}
