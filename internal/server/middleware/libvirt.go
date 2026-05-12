package middleware

import (
	"net/http"

	"go.uber.org/zap"
)

func LibvirtMiddleware(check func() bool, logger *zap.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if !check() {
				logger.Panic("libvirt connection unavailable", zap.String("path", r.URL.Path))
			}
			next.ServeHTTP(w, r)
		})
	}
}
