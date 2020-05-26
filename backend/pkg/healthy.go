package server

import (
	"context"
	"log"
	"net/http"
	"time"
)

type healthy interface {
	healthy(context.Context) error
}

func HealthHandler(h healthy) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), 1*time.Second)
		defer cancel()
		err := h.healthy(ctx)
		if err != nil {
			log.Printf("healthcheck: %s", err)
			w.WriteHeader(http.StatusServiceUnavailable)
			return
		}
		w.WriteHeader(http.StatusOK)
	}
}
