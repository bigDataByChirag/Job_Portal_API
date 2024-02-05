package middleware

import (
	"net/http"
	"time"

	"github.com/rs/zerolog/log"
)

// ResponseRecorder is a custom implementation of http.ResponseWriter that records the HTTP status code.
type ResponseRecorder struct {
	http.ResponseWriter
	StatusCode int
}

// WriteHeader overrides the WriteHeader method of http.ResponseWriter to record the HTTP status code.
func (rec *ResponseRecorder) WriteHeader(statusCode int) {
	rec.StatusCode = statusCode
	rec.ResponseWriter.WriteHeader(statusCode)
}

// HttpLogger is a middleware that logs information about incoming HTTP requests, including method, URL, status code, and duration.
// It wraps an existing http.Handler and returns a new http.Handler.
func HttpLogger(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Record the start time of the request processing.
		startTime := time.Now()
		// Create a custom ResponseRecorder to capture the status code.
		rec := &ResponseRecorder{
			ResponseWriter: w,
			// Initialize StatusCode to http.StatusOK if needed.
			// StatusCode: http.StatusOK,
		}
		// Serve the HTTP request using the wrapped handler.
		handler.ServeHTTP(rec, r)
		// Calculate the duration of the request processing.
		duration := time.Since(startTime)
		// Create a logger with request information.
		logger := log.Info()
		logger.
			Str("protocol", "http").
			Str("method", r.Method).
			Stringer("url", r.URL).
			Int("status", rec.StatusCode).
			Str("status_text", http.StatusText(rec.StatusCode)).
			Dur("duration", duration).
			Msg("Received an HTTP Request")
	})
}
