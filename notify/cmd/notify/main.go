package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go/config"
	cli "github.com/urfave/cli/v2"

	//	"gitlab.com/davedamoon/dinghy/backend/pkg/middleware"
	notify "gitlab.com/davedamoon/dinghy/notify/pkg"
	"gitlab.com/davedamoon/dinghy/notify/pkg/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	app := &cli.App{
		Name:  "notify",
		Usage: "Propagate calls to a http hook towards GRPC clients.",
		Flags: []cli.Flag{
			&cli.StringFlag{Name: "addr", Value: ":4080", Usage: "Address for server."},
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

	grpcServer := grpc.NewServer()
	pb.RegisterGreeterServer(grpcServer, &notify.GRPCServer{})
	reflection.Register(grpcServer)

	httpSrv := notify.NewServer()
	//svcHandler := middleware.RequestID(rand.Int63, httpServer)
	//svcHandler = middleware.InitTraceContext(svcHandler)
	//svcHandler = middleware.InstrumentHttpHandler(svcHandler)
	//svcHandler = middleware.Timeout(29*time.Second, svcHandler)

	svcServer := httpServer(grpcHTTPSwitch(grpcServer, httpSrv), c.String("addr"))

	log.Println("starting server")

	go mustListenAndServe(svcServer)

	log.Println("running")

	awaitShutdown()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	err = svcServer.Shutdown(ctx)
	if err != nil {
		return fmt.Errorf("shutdown server: %v", err)
	}

	grpcServer.Stop()

	return nil
}

func grpcHTTPSwitch(g *grpc.Server, h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("http version: %d, ct: %v", r.ProtoMajor, r.Header)
		if r.ProtoMajor == 2 { //&& strings.Contains(r.Header.Get("Content-Type"), "application/grpc") {
			g.ServeHTTP(w, r)
			return
		}

		h.ServeHTTP(w, r)
	})
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
	defer close(stop)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop
}
