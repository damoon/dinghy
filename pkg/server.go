package server

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type server struct {
	publicServer *http.Server
	adminServer  *http.Server
}

func NewServer(publicAddr, adminAddr string, publicHandler, healthHandler http.Handler) server {
	serviceMux := http.NewServeMux()
	serviceMux.Handle("/", http.TimeoutHandler(publicHandler, 30*time.Second, ""))
	adminMux := http.NewServeMux()
	adminMux.Handle("/healthz", healthHandler)
	adminMux.Handle("/metrics", promhttp.Handler())

	publicServer := &http.Server{
		Addr:         publicAddr,
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

	return server{
		publicServer: publicServer,
		adminServer:  adminServer,
	}
}

func (s server) Run(shutdown <-chan os.Signal) {
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
	err := s.publicServer.Shutdown(context.Background())
	if err != nil {
		log.Fatalf("server shutdown failed: %s\n", err)
	}
	err = s.adminServer.Shutdown(context.Background())
	if err != nil {
		log.Fatalf("server shutdown failed: %s\n", err)
	}
}
