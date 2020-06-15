package dinghy

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	//	"github.com/coreos/etcd/clientv3"
)

var notify = make(chan interface{})

func (s *ServiceServer) post(w http.ResponseWriter, r *http.Request) {
	if isMinioUpdate(r) {
		notifyWebsocket(w, r)
		return
	}
}

func isMinioUpdate(r *http.Request) bool {
	return r.Header.Get("Authorization") == "Bearer auth_token_value"
}

type MinioNotification struct {
	Key       string
	EventName string
}

func notifyWebsocket(w http.ResponseWriter, r *http.Request) {
	bodyBuffer, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("notifyWebsocket: parse json: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	notification := &MinioNotification{}
	json.Unmarshal(bodyBuffer, &notification)

	log.Println(notification.EventName)
	if notification.EventName != "s3:ObjectRemoved:Delete" && notification.EventName != "s3:ObjectCreated:Put" && notification.EventName != "s3:ObjectCreated:Copy" {
		return
	}

	//	log.Println(notification.EventName)
	//log.Println(string(bodyBuffer))
	notify <- nil
	//	notify()
}

//func notify(path, event string) {
//	// TODO set keep alive checks & timeouts
//	cli, err := clientv3.New(clientv3.Config{
//		Endpoints:   []string{"etcd:2379"},
//		DialTimeout: 1 * time.Second,
//	})
//	if err != nil {
//		// handle error!
//	}
//	defer cli.Close()
//
//	_, err := cli.Put(context.Backgroun(), path, event)
//	if err != nil {
//		// handle error!
//	}
//}
