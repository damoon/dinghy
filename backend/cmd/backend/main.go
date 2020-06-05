package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strings"
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
			&cli.StringFlag{Name: "service-addr", Value: ":8080", Usage: "Address for user service."},
			&cli.StringFlag{Name: "admin-addr", Value: ":9090", Usage: "Address for administration service."},
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
	log.Println("connect to minio")

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

	adm := dinghy.NewAdminServer()
	adm.Storage = storage
	admServer := httpServer(adm, c.String("admin-addr"))

	svc := dinghy.NewServiceServer()
	svc.Storage = storage
	svc.FrontendURL = c.String("frontend-url")
	svcServer := httpServer(svc, c.String("service-addr"))

	log.Println("start admin server")

	go mustListenAndServe(admServer)

	log.Println("start service server")

	go mustListenAndServe(svcServer)

	log.Println("fully started")

	awaitShutdown()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	err = shutdown(ctx, svcServer)
	if err != nil {
		return fmt.Errorf("shutdown service server: %v", err)
	}

	err = shutdown(ctx, admServer)
	if err != nil {
		return fmt.Errorf("shutdown admin server: %v", err)
	}

	return nil
}

func setupMinio(endpoint, accessKey, secretPath string, useSSL bool, region, bucket string) (*minio.Client, error) {
	secretKeyBytes, err := ioutil.ReadFile(secretPath)
	if err != nil {
		return nil, fmt.Errorf("reading secret access key from %s: %v", secretPath, err)
	}

	secretKey := strings.TrimSpace(string(secretKeyBytes))

	clnt, err := minio.NewWithRegion(endpoint, accessKey, secretKey, useSSL, region)
	if err != nil {
		return nil, err
	}

	//	http.DefaultTransport.ResponseHeaderTimeout = 10 * time.Second

	clnt.SetCustomTransport(&http.Transport{
		Dial: (&net.Dialer{
			Timeout: 5 * time.Second,
		}).Dial,
		TLSHandshakeTimeout: 5 * time.Second,
	})

	exists, err := clnt.BucketExists(bucket)
	if err != nil {
		return nil, fmt.Errorf("look up bucket %s: %v", bucket, err)
	}

	if !exists {
		return nil, fmt.Errorf("bucket %s does not exist", bucket)
	}

	return clnt, nil
}

func httpServer(h http.Handler, addr string) *http.Server {
	httpServer := &http.Server{
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
	httpServer.Addr = addr
	httpServer.Handler = h

	return httpServer
}

func mustListenAndServe(srv *http.Server) {
	err := srv.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		log.Fatal(err)
	}
}

func awaitShutdown() {
	stop := make(chan os.Signal, 2)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop
}

func shutdown(ctx context.Context, srv *http.Server) error {
	err := srv.Shutdown(ctx)
	if err != nil {
		return err
	}

	return nil
}
