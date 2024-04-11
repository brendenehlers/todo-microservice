package http

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/brendenehlers/todo-microservice/domain"
)

type Middleware func(handler http.Handler) http.Handler

func RequestLogger(log domain.Logger, handler http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		defer func() {
			log.Info(fmt.Sprintf("[%s] %s %v", r.Method, r.URL.Path, time.Since(start)))
		}()

		handler.ServeHTTP(w, r)
	}

	return http.HandlerFunc(fn)
}

// TODO finish this
func RequestTimeout(timeout time.Duration, handler http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), timeout)
		defer cancel()

		go handler.ServeHTTP(w, r)

		<-ctx.Done()
		fmt.Println("context reacted")
		cancel()
	}

	return http.HandlerFunc(fn)
}
