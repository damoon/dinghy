package main

import (
	"flag"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	server "webdav-to-s3/pkg"

	"github.com/minio/minio-go"
)

func main() {
	serviceAddr := flag.String("service-address", ":8080", "service server address, ':8080'")
	adminAddr := flag.String("admin-address", ":8081", "admin server address, ':8081'")
	endpoint := flag.String("endpoint", "http://127.0.0.1:9000", "s3 endpoint")
	accessKeyID := flag.String("accessKeyID", "", "s3 accessKeyID")
	secretAccessKey := flag.String("secretAccessKey", "", "s3 secretAccessKey")
	useSSL := flag.Bool("useSSL", true, "s3 uses https")
	bucket := flag.String("bucket", "webdav", "s3 bucket name")
	location := flag.String("location", "us-east-1", "s3 bucket location")
	redirectURL := flag.String("redirectURL", "http://127.0.0.1:9000", "url to redirect to instead of 404 (minio)")
	lightWeight := flag.Bool("light", true, "only support GET and PUT via redirects")

	flag.Parse()

	log.Printf("serviceAddr: %s\n", *serviceAddr)
	log.Printf("adminAddr: %s\n", *adminAddr)
	log.Printf("endpoint: %s\n", *endpoint)
	log.Printf("accessKeyID: %s\n", *accessKeyID)
	log.Printf("secretAccessKey: %s\n", *secretAccessKey)
	log.Printf("useSSL: %v\n", *useSSL)
	log.Printf("bucket: %s\n", *bucket)
	log.Printf("location: %s\n", *location)
	log.Printf("redirectURL: %s\n", *redirectURL)
	log.Printf("lightWeight: %v\n", *lightWeight)

	transport := &http.Transport{
		Dial: (&net.Dialer{
			Timeout: 5 * time.Second,
		}).Dial,
		TLSHandshakeTimeout: 5 * time.Second,
	}

	minioClient, err := minio.New(*endpoint, *accessKeyID, *secretAccessKey, *useSSL)
	if err != nil {
		log.Fatalln(err)
	}
	minioClient.SetCustomTransport(transport)

	storage := server.NewMinioStorage(minioClient, *bucket, *location)
	go storage.EnsureBucket()

	healthHandler := server.HealthHandler(storage)
	serviceHandler := server.NewPresignHandler(storage, *redirectURL)
	if !*lightWeight {
		serviceHandler = server.NewForwardHandler(storage)
	}

	// run server until exit signal
	stop := make(chan os.Signal, 2)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	s := server.NewServer(*serviceAddr, *adminAddr, serviceHandler, healthHandler)
	s.Run(stop)
}
