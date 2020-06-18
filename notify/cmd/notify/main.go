package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go/config"
	cli "github.com/urfave/cli/v2"

	//	"gitlab.com/davedamoon/dinghy/backend/pkg/middleware"
	notify "gitlab.com/davedamoon/dinghy/notify/pkg"
	"gitlab.com/davedamoon/dinghy/notify/pkg/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"
)

func main() {
	app := &cli.App{
		Name:  "notify",
		Usage: "Propagate calls to a http hook towards GRPC clients.",
		Flags: []cli.Flag{
			&cli.StringFlag{Name: "http", Value: ":8080", Usage: "Address for server."},
			&cli.StringFlag{Name: "grpc", Value: ":50051", Usage: "Address for server."},
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

	jaeger, err := setupJaeger()
	if err != nil {
		return fmt.Errorf("setup minio s3 client: %v", err)
	}
	defer jaeger.Close()

	m := sync.Mutex{}
	cond := sync.NewCond(&m)

	grpcServce := &notify.GRPCServer{}
	grpcServce.C = cond

	httpSrv := notify.NewServer()
	httpSrv.C = cond
	//svcHandler := middleware.RequestID(rand.Int63, httpServer)
	//svcHandler = middleware.InitTraceContext(svcHandler)
	//svcHandler = middleware.InstrumentHttpHandler(svcHandler)
	//svcHandler = middleware.Timeout(29*time.Second, svcHandler)

	grpcS := grpc.NewServer()
	pb.RegisterNotifierServer(grpcS, grpcServce)
	grpc_health_v1.RegisterHealthServer(grpcS, health.NewServer())
	reflection.Register(grpcS)

	httpS := httpServer(httpSrv, c.String("http"))

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
