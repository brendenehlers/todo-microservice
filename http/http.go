package http

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/brendenehlers/todo-microservice/domain"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
)

const (
	REQUEST_TIMEOUT = time.Millisecond * 200
	DEFAULT_ADDRESS = ":8080"
)

var (
	ErrRequestTimedOut = fmt.Errorf("request timed out")
	ErrInvalidRepo     = fmt.Errorf("invalid todo repository")
	ErrInvalidLogger   = fmt.Errorf("invalid logger")
	ErrNoPathValue     = fmt.Errorf("no path value found")
)

type HTTPServerConfig struct {
	Addr string
	Repo domain.TodoRepository
	Ctx  context.Context
	Log  domain.Logger
}

func CreateHTTPServer(config *HTTPServerConfig) (*HttpServer, error) {
	if config.Repo == nil {
		return nil, ErrInvalidRepo
	}
	if config.Log == nil {
		return nil, ErrInvalidLogger
	}

	if config.Ctx == nil {
		config.Ctx = context.Background()
	}
	if config.Addr == "" {
		config.Addr = DEFAULT_ADDRESS
	}

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	repoAdapter := newAdapter(config.Repo)

	api := newAPI(
		repoAdapter,
		config.Log,
	)
	r.Mount("/", api.handler())

	server := &HttpServer{
		Server: http.Server{
			Addr:    config.Addr,
			Handler: r,
			BaseContext: func(l net.Listener) context.Context {
				return config.Ctx
			},
		},
		log: config.Log,
	}

	return server, nil
}

type HttpServer struct {
	http.Server
	log domain.Logger
}

func (s *HttpServer) Run() {
	s.log.Info(fmt.Sprintf("Server running on %s", s.Addr))
	s.ListenAndServe()
}

func (s *HttpServer) Stop(ctx context.Context) {
	s.Shutdown(ctx)
}
