package middleware

import (
	"net/http"
	"sync/atomic"
)

func RateLimiter(limit uint32, handler http.Handler) http.Handler {
	var counter int64
	limitInt64 := int64(limit)
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if atomic.LoadInt64(&counter) >= limitInt64 {
			http.Error(w, http.StatusText(http.StatusTooManyRequests), http.StatusTooManyRequests)
		} else {
			atomic.AddInt64(&counter, 1)
			handler.ServeHTTP(w, r)
			atomic.AddInt64(&counter, -1)
		}
	})
}
