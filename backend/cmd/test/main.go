package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"time"

	minio "github.com/minio/minio-go"
)

func main() {
	transport := &http.Transport{
		Dial: (&net.Dialer{
			Timeout: 5 * time.Second,
		}).Dial,
		TLSHandshakeTimeout: 5 * time.Second,
	}

	minioClient, err := minio.New("minio:9000", "minio", "minio123", false)
	if err != nil {
		log.Fatalf("set up minio client: %v", err)
	}
	minioClient.SetCustomTransport(transport)

	// Create a done channel to control 'ListObjectsV2' go routine.
	doneCh := make(chan struct{})

	// Indicate to our routine to exit cleanly upon return.
	defer close(doneCh)

	isRecursive := true
	prefix := ""
	objectCh := minioClient.ListObjectsV2("dinghy", prefix, isRecursive, doneCh)
	for object := range objectCh {
		if object.Err != nil {
			fmt.Println(object.Err)
			return
		}
		fmt.Println(object)
	}
}
