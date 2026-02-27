package middleware

import (
	"net/http"

	"go.uber.org/zap"
)

func LibvirtMiddleware(check func() bool, logger *zap.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if !check() {
				// find a way for graceful shutdown of server when libvirt connection is unavailable
				logger.Panic("libvirt connection unavailable", zap.String("path", r.URL.Path))
				http.Error(w, "service unavailable", http.StatusServiceUnavailable)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
