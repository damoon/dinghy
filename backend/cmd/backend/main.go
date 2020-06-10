package main

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go/config"
	cli "github.com/urfave/cli/v2"
	dinghy "gitlab.com/davedamoon/dinghy/backend/pkg"
	"gitlab.com/davedamoon/dinghy/backend/pkg/middleware"
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
		log.Println(err)
		os.Exit(1)
	}

	log.Println("shutdown complete")
}

func run(c *cli.Context) error {
	log.Println("set up tracing")

	err, jaeger := setupJaeger()
	if err != nil {
		return fmt.Errorf("setup minio s3 client: %v", err)
	}
	defer jaeger.Close()

	log.Println("set up storage")

	s3Client, err := setupMinio(
		c.String("s3-endpoint"),
		c.String("s3-access-key"),
		c.String("s3-secret-access-key-file"),
		c.Bool("s3-ssl"),
		c.String("s3-location"),
		c.String("s3-bucket"))
	if err != nil {
		return fmt.Errorf("setup minio s3 client: %v", err)
	}

	storage := dinghy.Storage{
		Client: s3Client,
		Bucket: c.String("s3-bucket"),
	}

	adm := dinghy.NewAdminServer()
	adm.Storage = storage
	admHandler := middleware.RequestID(rand.Int63, adm)
	admHandler = middleware.InitTraceContext(admHandler)
	//admHandler = dinghy.InstrumentHttpHandler(admHandler) // reduce noise
	admHandler = middleware.Timeout(900*time.Millisecond, admHandler)
	admServer := httpServer(admHandler, c.String("admin-addr"))

	svc := dinghy.NewServiceServer()
	svc.Storage = storage
	svc.FrontendURL = c.String("frontend-url")
	svcHandler := middleware.CORS(c.String("frontend-url"), svc)
	svcHandler = middleware.RequestID(rand.Int63, svcHandler)
	svcHandler = middleware.InitTraceContext(svcHandler)
	svcHandler = middleware.InstrumentHttpHandler(svcHandler)
	svcHandler = middleware.Timeout(29*time.Second, svcHandler)
	svcServer := httpServer(svcHandler, c.String("service-addr"))

	log.Println("starting admin server")

	go mustListenAndServe(admServer)

	log.Println("starting service server")

	go mustListenAndServe(svcServer)

	log.Println("running")

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

func setupMinio(endpoint, accessKey, secretPath string, useSSL bool, region, bucket string) (*s3.S3, error) {
	secretKeyBytes, err := ioutil.ReadFile(secretPath)
	if err != nil {
		return nil, fmt.Errorf("reading secret access key from %s: %v", secretPath, err)
	}

	secretKey := strings.TrimSpace(string(secretKeyBytes))

	endpointProtocol := "http"
	if useSSL {
		endpointProtocol = "https"
	}

	s3Config := &aws.Config{
		Credentials:      credentials.NewStaticCredentials(accessKey, secretKey, ""),
		Endpoint:         aws.String(fmt.Sprintf("%s://%s", endpointProtocol, endpoint)),
		Region:           aws.String(region),
		DisableSSL:       aws.Bool(!useSSL),
		S3ForcePathStyle: aws.Bool(true),
	}

	newSession, err := session.NewSession(s3Config)
	if err != nil {
		return nil, fmt.Errorf("set up aws session: %v", err)
	}

	s3Client := s3.New(newSession)

	return s3Client, nil
}

func setupJaeger() (error, io.Closer) {
	cfg, err := config.FromEnv()
	if err != nil {
		return fmt.Errorf("load config from Env Vars: %v", err), nil
	}

	tracer, closer, err := cfg.NewTracer()
	if err != nil {
		return err, nil
	}

	opentracing.SetGlobalTracer(tracer)

	return nil, closer
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
