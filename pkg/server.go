package server

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/minio/minio-go"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func RunServer(mc *minio.Client, bucket, adminAddr, serviceAddr, redirectURL string, minimalizm bool) {
	serviceMux := http.NewServeMux()
	var handler http.HandlerFunc
	switch minimalizm {
	case true:
		handler = presign(mc, bucket, redirectURL)
	case false:
		handler = emulate(mc, bucket, redirectURL)
	}
	serviceMux.Handle("/", http.TimeoutHandler(handler, 30*time.Second, ""))
	adminMux := http.NewServeMux()
	adminMux.Handle("/healthz", http.TimeoutHandler(newHealth(mc, bucket), 9*time.Second, ""))
	adminMux.Handle("/metrics", promhttp.Handler())

	serviceServer := &http.Server{
		Addr:         serviceAddr,
		Handler:      serviceMux,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
	}
	adminServer := &http.Server{
		Addr:         adminAddr,
		Handler:      adminMux,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
	}

	go func() {
		err := serviceServer.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()
	go func() {
		err := adminServer.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()

	// wait for an exit signal
	stop := make(chan os.Signal, 2)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	<-stop
	err := serviceServer.Shutdown(context.Background())
	if err != nil {
		log.Fatalf("server shutdown failed: %s\n", err)
	}
	err = adminServer.Shutdown(context.Background())
	if err != nil {
		log.Fatalf("server shutdown failed: %s\n", err)
	}
}
