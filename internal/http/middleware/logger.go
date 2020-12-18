package middleware

import (
	"net/http"

	"github.com/vlad161/affise_test_task/internal/logger"
)

func Logger(log logger.Logger, handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		ctx = logger.WithContext(ctx, log)
		r = r.WithContext(ctx)
		handler.ServeHTTP(w, r)
	})
}
