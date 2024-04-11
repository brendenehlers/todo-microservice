package http

import (
	"fmt"
	"net/http"
	"runtime/debug"
	"time"

	"github.com/brendenehlers/todo-microservice/domain"
)

type middleware struct {
	log domain.Logger
}

func (mw *middleware) requestLogger(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		defer func() {
			mw.log.Info(fmt.Sprintf("[%s] %s %v", r.Method, r.URL.Path, time.Since(start)))
		}()

		next.ServeHTTP(w, r)
	}

	return http.HandlerFunc(fn)
}

func (mw *middleware) recover(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rvr := recover(); rvr != nil {
				if rvr == http.ErrAbortHandler {
					// continue panicking if ErrAbortHandler
					panic(rvr)
				}

				mw.log.Error(fmt.Sprintf("Panic: %v", rvr))
				debug.PrintStack()
			}
		}()

		next.ServeHTTP(w, r)
	}

	return http.HandlerFunc(fn)
}

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
