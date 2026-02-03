package middleware

import (
	"log"
	"net/http"
	"time"
)

type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		start := time.Now()

		// Wrap ResponseWriter to capture status code
		wrapped := &responseWriter{ResponseWriter: writer, statusCode: http.StatusOK}

		next.ServeHTTP(wrapped, request)

		log.Printf(
			"%s %s %d %s",
			request.Method,
			request.RequestURI,
			wrapped.statusCode,
			time.Since(start),
		)
	})
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}
