package middleware

import (
	"net/http"
)

func CORS(domain string, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", domain)
		next.ServeHTTP(w, r)
	})
}
