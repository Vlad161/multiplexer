package middleware

import (
	"net/http"
)

func RateLimiter(limit uint32, handler http.Handler) http.Handler {
	ch := make(chan struct{}, limit)

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		select {
		case ch <- struct{}{}:
			handler.ServeHTTP(w, r)
			<-ch
		default:
			http.Error(w, http.StatusText(http.StatusTooManyRequests), http.StatusTooManyRequests)
		}
	})
}
