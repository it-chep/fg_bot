package middleware

import (
	"net/http"

	"fg_bot/internal/pkg/logger"
)

func LoggerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		ctx = logger.ContextWithLogger(ctx, logger.New())
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}
