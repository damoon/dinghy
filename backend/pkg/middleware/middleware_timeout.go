package middleware

import (
	"context"
	"net/http"
	"time"
)

func Timeout(timeout time.Duration, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), timeout)
		defer cancel()

		done := make(chan bool, 1)

		go func() {
			next.ServeHTTP(w, r.WithContext(ctx))
			done <- true
			close(done)
		}()

		select {
		case <-ctx.Done():
			if ctx.Err() == context.DeadlineExceeded {
				w.WriteHeader(http.StatusGatewayTimeout)
			}
		case <-done:
		}
	})
}
