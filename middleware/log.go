package middleware

import (
	"log"
	"net/http"
	"time"
)

// responseRecorder - записывает статус ответа.
type responseRecorder struct {
	http.ResponseWriter
	statusCode int
}

func (r *responseRecorder) WriteHeader(statusCode int) {
	r.statusCode = statusCode
	r.ResponseWriter.WriteHeader(statusCode)
}

// LoggingMiddleware - мидлварь для логирования запросов.
func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rr := &responseRecorder{w, http.StatusOK}
		start := time.Now()

		next.ServeHTTP(rr, r)

		log.Printf("[%s] %s %s %d %s | %v",
			time.Now().Format(time.RFC1123),
			r.Method,
			r.URL.Path,
			rr.statusCode,
			r.UserAgent(),
			time.Since(start),
		)
	})
}
