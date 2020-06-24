package main

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"math/rand"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	"time"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	otgrpc "github.com/opentracing-contrib/go-grpc"
	opentracing "github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go/config"
	cli "github.com/urfave/cli/v2"
	notify "gitlab.com/davedamoon/dinghy/notify/pkg"
	"gitlab.com/davedamoon/dinghy/notify/pkg/middleware"
	"gitlab.com/davedamoon/dinghy/notify/pkg/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"
)

var (
	gitHash string
	gitRef  string
)

func main() {
	app := &cli.App{
		Name:  "notify",
		Usage: "Propagate notifications from a http hook towards GRPC clients.",
		Commands: []*cli.Command{
			{
				Name:  "server",
				Usage: "Start the server.",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "http", Value: ":8080", Usage: "Address for server."},
					&cli.StringFlag{Name: "grpc", Value: ":50051", Usage: "Address for server."},
					&cli.StringFlag{Name: "webhook-token-file", Required: true, Usage: "Path to webhook token file."},
				},
				Action: run,
			},
			{
				Name:  "version",
				Usage: "Show the version",
				Action: func(c *cli.Context) error {
					log.Printf("%s - %s", gitRef, gitHash)
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
	log.Printf("git hash: %v", gitHash)
	log.Printf("git ref: %v", gitRef)

	tokenBytes, err := ioutil.ReadFile(c.String("webhook-token-file"))
	if err != nil {
		return fmt.Errorf("reading webhook token from %s: %v", c.String("webhook-token-file"), err)
	}

	token := strings.TrimSpace(string(tokenBytes))

	log.Println("set up metrics")

	middleware.InitMetrics(gitHash, gitRef)

	log.Println("set up tracing")

	jaeger, err := setupJaeger()
	if err != nil {
		return fmt.Errorf("setup tracing: %v", err)
	}
	defer jaeger.Close()

	log.Println("set up servers")

	m := sync.Mutex{}
	cond := sync.NewCond(&m)

	grpcServce := &notify.GRPCServer{}
	grpcServce.C = cond

	httpSrv := notify.NewServer()
	httpSrv.BearerToken = token
	httpSrv.C = cond
	svcHandler := middleware.RequestID(rand.Int63, httpSrv)
	svcHandler = middleware.InitTraceContext(svcHandler)
	svcHandler = middleware.InstrumentHttpHandler(svcHandler)
	svcHandler = middleware.Timeout(29*time.Second, svcHandler)

	tracer := opentracing.GlobalTracer()
	grpcS := grpc.NewServer(
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
			grpc_prometheus.UnaryServerInterceptor,
			otgrpc.OpenTracingServerInterceptor(tracer))),
		grpc.StreamInterceptor(grpc_middleware.ChainStreamServer(
			grpc_prometheus.StreamServerInterceptor,
			otgrpc.OpenTracingStreamServerInterceptor(tracer))))
	pb.RegisterNotifierServer(grpcS, grpcServce)
	grpc_prometheus.Register(grpcS)
	grpc_health_v1.RegisterHealthServer(grpcS, health.NewServer())
	reflection.Register(grpcS)

	httpS := httpServer(svcHandler, c.String("http"))

	log.Println("starting grpc server")

	go mustListenAndServeGRPC(grpcS, c.String("grpc"))

	log.Println("starting http server")

	go mustListenAndServeHTTP(httpS)

	log.Println("running")

	awaitShutdown()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	err = httpS.Shutdown(ctx)
	if err != nil {
		return fmt.Errorf("shutdown server: %v", err)
	}

	grpcS.Stop()

	log.Println("shutdown complete")

	return nil
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

func mustListenAndServeHTTP(srv *http.Server) {
	err := srv.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		log.Fatal(err)
	}
}

func mustListenAndServeGRPC(s *grpc.Server, addr string) {
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

func awaitShutdown() {
	stop := make(chan os.Signal, 2)
	defer close(stop)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop
}
