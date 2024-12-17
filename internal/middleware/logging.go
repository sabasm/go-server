package middleware

import (
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

func LoggingMiddleware(logger *zap.Logger) mux.MiddlewareFunc {
	httpLogger := logger.Named("http")

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			next.ServeHTTP(w, r)
			duration := time.Since(start)

			// Only log if it's not a GET request or takes longer than expected
			if r.Method != http.MethodGet || duration > time.Millisecond {
				httpLogger.Info("request",
					zap.String("method", r.Method),
					zap.String("path", r.URL.Path),
					zap.Duration("took", duration),
					zap.String("from", r.RemoteAddr),
					zap.String("agent", r.UserAgent()),
				)
			}
		})
	}
}
