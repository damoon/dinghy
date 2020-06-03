package dinghy

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type Server struct {
	publicServer *http.Server
	adminServer  *http.Server
}

func NewServer(publicAddr, adminAddr string, publicHandler, healthHandler http.Handler) Server {
	serviceMux := http.NewServeMux()
	svc := &svc{}
	serviceMux.Handle("/", http.TimeoutHandler(svc, 30*time.Second, ""))
	adminMux := http.NewServeMux()
	adminMux.Handle("/healthz", healthHandler)
	adminMux.Handle("/metrics", promhttp.Handler())

	publicServer := &http.Server{
		Addr:         publicAddr,
		Handler:      serviceMux,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
	adminServer := &http.Server{
		Addr:         adminAddr,
		Handler:      adminMux,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	return Server{
		publicServer: publicServer,
		adminServer:  adminServer,
	}
}

func (s Server) Run(shutdown <-chan os.Signal) {
	go func() {
		err := s.publicServer.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()
	go func() {
		err := s.adminServer.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()

	<-shutdown

	ctx, cancel := context.WithTimeout(context.Background(), 35*time.Second)
	defer cancel()
	err := s.publicServer.Shutdown(ctx)
	if err != nil {
		log.Fatalf("server shutdown failed: %s\n", err)
	}

	ctx, cancel = context.WithTimeout(context.Background(), 35*time.Second)
	defer cancel()
	err = s.adminServer.Shutdown(ctx)
	if err != nil {
		log.Fatalf("server shutdown failed: %s\n", err)
	}
}
