package middleware

import (
	"log"
	"net/http"
	"time"
)

func JSONContentTypeMiddleware(next http.Handler) http.Handler{
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-type", "application/json")
		next.ServeHTTP(w, r)
	})
}

func LoggerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		next.ServeHTTP(w, r)

		log.Printf("[%s] %s - Tamamlanma süresi: %v", r.Method, r.URL.Path, time.Since(start))
	})
}