package log

import (
	"net/http"
	"os"
	"time"

	"log/slog"
)

func LogMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelInfo,
		})
		logger := slog.New(handler)

		lrw := &loggingResponseWriter{w, http.StatusOK}

		next.ServeHTTP(lrw, r)

		latency := time.Since(start)
		logger.Info(
			"request_details",
			slog.String("method", r.Method),
			slog.Int("status", lrw.status),
			slog.String("path", r.URL.Path),
			slog.Duration("latency", latency),
			slog.String("client_ip", r.RemoteAddr),
			slog.String("timestamp", time.Now().Format("2006/01/02 - 15:04:05")),
		)
	})
}

type loggingResponseWriter struct {
	http.ResponseWriter
	status int
}

func (lrw *loggingResponseWriter) WriteHeader(statusCode int) {
	lrw.status = statusCode
	lrw.ResponseWriter.WriteHeader(statusCode)
}

func (lrw *loggingResponseWriter) Write(b []byte) (int, error) {
	return lrw.ResponseWriter.Write(b)
}
