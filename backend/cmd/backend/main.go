package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	minio "github.com/minio/minio-go"
	"github.com/urfave/cli/v2"
	dinghy "gitlab.com/davedamoon/dinghy/backend/pkg"
)

func main() {
	app := &cli.App{
		Name:  "boom",
		Usage: "make an explosive entrance",
		Flags: []cli.Flag{
			&cli.StringFlag{Name: "user-address", Value: ":8080", Usage: "Address for user service."},
			&cli.StringFlag{Name: "admin-address", Value: ":9090", Usage: "Address for administration service."},
			&cli.StringFlag{Name: "s3-endpoint", Required: true, Usage: "s3 endpoint."},
			&cli.StringFlag{Name: "s3-access-key", Required: true, Usage: "s3 access key."},
			&cli.StringFlag{Name: "s3-secret-access-key-file", Required: true, Usage: "Path to s3 secret access key."},
			&cli.BoolFlag{Name: "s3-ssl", Value: true, Usage: "s3 uses SSL."},
			&cli.StringFlag{Name: "s3-location", Value: "us-east-1", Usage: "s3 bucket location."},
			&cli.StringFlag{Name: "s3-bucket", Required: true, Usage: "s3 bucket name."},
			&cli.StringFlag{Name: "frontend-url", Required: true, Usage: "Frontend domain for CORS and redirects."},
		},
		Action: run,
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func run(c *cli.Context) error {

	s3Clnt, err := setupMinio(
		c.String("s3-endpoint"),
		c.String("s3-access-key"),
		c.String("s3-secret-access-key-file"),
		c.Bool("s3-ssl"),
		c.String("s3-location"),
		c.String("s3-bucket"))
	if err != nil {
		return fmt.Errorf("setup minio s3 client: %v", err)
	}

	storage := dinghy.NewMinioStorage(s3Clnt, c.String("s3-location"), c.String("s3-bucket"))
	go storage.EnsureBucket()

	healthHandler := dinghy.HealthHandler(storage)
	serviceHandler := dinghy.NewPresignHandler(storage, c.String("frontend-url"))
	serviceHandler = dinghy.NewForwardHandler(storage)

	// run server until exit signal
	stop := make(chan os.Signal, 2)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	s := dinghy.NewServer(c.String("user-service"), c.String("admin-service"), serviceHandler, healthHandler)
	s.Run(stop)

	return nil
}

func setupMinio(endpoint, accessKey, secretPath string, useSSL bool, region, bucket string) (*minio.Client, error) {
	secretAccessKey, err := ioutil.ReadFile(secretPath)
	if err != nil {
		return nil, fmt.Errorf("reading secret access key from %s: %v", secretPath, err)
	}

	clnt, err := minio.NewWithRegion(endpoint, accessKey, string(secretAccessKey), useSSL, region)
	if err != nil {
		return nil, err
	}

	exists, err := clnt.BucketExists(bucket)
	if err != nil {
		return nil, fmt.Errorf("look up bucket %s: %v", bucket, err)
	}

	if !exists {
		return nil, fmt.Errorf("bucket %s does not exist", bucket)
	}

	clnt.SetCustomTransport(&http.Transport{
		Dial: (&net.Dialer{
			Timeout: 5 * time.Second,
		}).Dial,
		TLSHandshakeTimeout: 5 * time.Second,
	})

	return clnt, nil
}
