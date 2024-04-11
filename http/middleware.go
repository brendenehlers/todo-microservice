package http

import (
	"net/http"
	"time"
)

// TODO finish this
func RequestTimeout(timeout time.Duration, handler http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		// ctx, cancel := context.WithTimeout(r.Context(), timeout)
		// defer cancel()

		// go handler.ServeHTTP(w, r)

		// <-ctx.Done()
		// fmt.Println("context reacted")
		// cancel()
	}

	return http.HandlerFunc(fn)
}
