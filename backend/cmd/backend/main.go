package main

import (
	"context"
	"crypto/tls"
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
	"github.com/gorilla/websocket"
	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	otgrpc "github.com/opentracing-contrib/go-grpc"
	opentracing "github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go/config"
	cli "github.com/urfave/cli/v2"
	dinghy "gitlab.com/davedamoon/dinghy/backend/pkg"
	"gitlab.com/davedamoon/dinghy/backend/pkg/middleware"
	"gitlab.com/davedamoon/dinghy/backend/pkg/pb"
	"google.golang.org/grpc"
)

var (
	gitHash string
	gitRef  string
)

func main() {
	app := &cli.App{
		Name:   "backend",
		Usage:  "Connect http to s3.",
		Action: run,
		Commands: []*cli.Command{
			{
				Name:  "server",
				Usage: "Start the server.",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "service-addr", Value: ":8080", Usage: "Address for user service."},
					&cli.StringFlag{Name: "admin-addr", Value: ":8090", Usage: "Address for administration service."},
					&cli.StringFlag{Name: "s3-endpoint", Required: true, Usage: "s3 endpoint."},
					&cli.StringFlag{Name: "s3-access-key-file", Required: true, Usage: "Path to s3 access key."},
					&cli.StringFlag{Name: "s3-secret-key-file", Required: true, Usage: "Path to s3 secret access key."},
					&cli.BoolFlag{Name: "s3-ssl", Value: true, Usage: "s3 uses SSL."},
					&cli.StringFlag{Name: "s3-location", Value: "us-east-1", Usage: "s3 bucket location."},
					&cli.StringFlag{Name: "s3-bucket", Required: true, Usage: "s3 bucket name."},
					&cli.StringFlag{Name: "frontend-url", Required: true, Usage: "Frontend domain for CORS and redirects."},
					&cli.StringFlag{Name: "notify-endpoint", Value: "notify:50051", Usage: "Notify service endpoint."},
					&cli.BoolFlag{Name: "tls-insecure-skip-verify", Value: false, Usage: "Disable TLS verification."},
				},
				Action: run,
			},
			{
				Name:  "version",
				Usage: "Show the version",
				Action: func(c *cli.Context) error {
					_, err := os.Stdout.WriteString(fmt.Sprintf("version: %s\ngit commit: %s", gitRef, gitHash))
					if err != nil {
						return err
					}

					return nil
				},
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
}

func run(c *cli.Context) error {
	log.Printf("version: %v", gitRef)
	log.Printf("git commit: %v", gitHash)

	if c.Bool("tls-insecure-skip-verify") {
		log.Println("insecure: disable TLS verification")

		http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	}

	log.Println("set up metrics")

	middleware.InitMetrics(gitHash, gitRef)

	log.Println("set up tracing")

	jaeger, err := setupJaeger()
	if err != nil {
		return fmt.Errorf("setup tracing: %v", err)
	}
	defer jaeger.Close()

	log.Println("set up storage")

	storage, err := setupMinioAdapter(
		c.String("s3-endpoint"),
		c.String("s3-access-key-file"),
		c.String("s3-secret-key-file"),
		c.Bool("s3-ssl"),
		c.String("s3-location"),
		c.String("s3-bucket"))
	if err != nil {
		return fmt.Errorf("setup minio s3 client: %v", err)
	}

	log.Println("set up notify client")

	nc, closeNotify, err := setupNotifyClient(c.String("notify-endpoint"))
	if err != nil {
		return fmt.Errorf("setup notify client: %v", err)
	}
	defer closeNotify.Close()

	log.Println("set up servers")

	adm := dinghy.NewAdminServer()
	adm.Storage = storage
	admHandler := middleware.RequestID(rand.Int63, adm)
	admHandler = middleware.InitTraceContext(admHandler)
	//admHandler = dinghy.InstrumentHttpHandler(admHandler) // reduce noise
	admHandler = middleware.Timeout(900*time.Millisecond, admHandler)
	admServer := httpServer(admHandler, c.String("admin-addr"))

	svc := &dinghy.ServiceServer{}
	svc.FrontendURL = c.String("frontend-url")
	svc.Storage = storage
	svc.Notify = nc
	svc.Upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			return r.Header.Get("Origin") == c.String("frontend-url")
		},
	}

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

	log.Println("shutdown complete")

	return nil
}

func setupNotifyClient(addr string) (*dinghy.NotifyAdapter, io.Closer, error) {
	tracer := opentracing.GlobalTracer()

	conn, err := grpc.Dial(
		addr,
		grpc.WithInsecure(),
		grpc.WithUnaryInterceptor(grpc_prometheus.UnaryClientInterceptor),
		grpc.WithUnaryInterceptor(otgrpc.OpenTracingClientInterceptor(tracer)),
		grpc.WithStreamInterceptor(grpc_prometheus.StreamClientInterceptor),
		grpc.WithStreamInterceptor(otgrpc.OpenTracingStreamClientInterceptor(tracer)))
	if err != nil {
		return nil, nil, fmt.Errorf("connect to %s: %v", addr, err)
	}

	return &dinghy.NotifyAdapter{
		NotifierClient: pb.NewNotifierClient(conn),
	}, conn, nil
}

func setupMinioAdapter(endpoint, accessKeyPath, secretKeyPath string,
	useSSL bool, region, bucket string) (*dinghy.MinioAdapter, error) {
	accessKeyBytes, err := ioutil.ReadFile(accessKeyPath)
	if err != nil {
		return nil, fmt.Errorf("reading secret access key from %s: %v", accessKeyPath, err)
	}

	secretKeyBytes, err := ioutil.ReadFile(secretKeyPath)
	if err != nil {
		return nil, fmt.Errorf("reading secret access key from %s: %v", secretKeyPath, err)
	}

	accessKey := strings.TrimSpace(string(accessKeyBytes))
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

	return &dinghy.MinioAdapter{
		Client: s3Client,
		Bucket: bucket,
	}, nil
}

func setupJaeger() (io.Closer, error) {
	cfg, err := config.FromEnv()
	if err != nil {
		return nil, fmt.Errorf("load config from Env Vars: %v", err)
	}

	tracer, closer, err := cfg.NewTracer()
	if err != nil {
		return nil, err
	}

	opentracing.SetGlobalTracer(tracer)

	return closer, nil
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
