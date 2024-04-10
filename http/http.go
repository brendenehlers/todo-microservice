package http

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/brendenehlers/todo-microservice"
)

const (
	REQUEST_TIMEOUT = time.Millisecond * 200
)

type HTTPServerConfig struct {
	Addr string
	Repo todo.TodoRepository
	Ctx  context.Context
	Log  todo.Logger
}

func CreateHTTPServer(config *HTTPServerConfig) (*HttpServer, error) {
	if config.Repo == nil {
		return nil, fmt.Errorf("config.todoRepo == nil")
	}
	if config.Log == nil {
		return nil, fmt.Errorf("config.log == nil")
	}

	if config.Ctx == nil {
		config.Ctx = context.Background()
	}
	if config.Addr == "" {
		config.Addr = ":8080"
	}

	handler := http.NewServeMux()
	server := &HttpServer{
		Server: http.Server{
			Addr:    config.Addr,
			Handler: RequestLogger(config.Log, handler),
			BaseContext: func(l net.Listener) context.Context {
				return config.Ctx
			},
		},
		repo: config.Repo,
		log:  config.Log,
	}

	handler.HandleFunc("POST /todo", server.handleCreateTodo)
	handler.HandleFunc("GET /todo/{todoId}", server.handleGetTodo)
	handler.HandleFunc("PUT /todo/{todoId}", server.handleUpdateTodo)
	handler.HandleFunc("DELETE /todo/{todoId}", server.handleDeleteTodo)

	return server, nil
}

type HttpServer struct {
	http.Server
	repo todo.TodoRepository
	log  todo.Logger
}

type serverResponse struct {
	Message string      `json:"message,omitempty"`
	Value   interface{} `json:"value,omitempty"`
	Error   string      `json:"error,omitempty"`
}

func (s *HttpServer) Run() {
	s.log.Info(fmt.Sprintf("Server running on %s", s.Addr))
	s.ListenAndServe()
}

func (*HttpServer) Stop() {
	panic("not implemented")
}

func (s *HttpServer) handleCreateTodo(w http.ResponseWriter, r *http.Request) {
	type response struct {
		val *todo.Todo
		err error
	}

	ctx, cancel := context.WithTimeout(r.Context(), REQUEST_TIMEOUT)
	defer cancel()
	respch := make(chan response)

	var newTodo todo.NewTodo
	json.NewDecoder(r.Body).Decode(&newTodo)
	defer r.Body.Close()

	go func() {
		todo, err := s.repo.CreateTodo(&newTodo)
		respch <- response{
			val: todo,
			err: err,
		}
	}()

	select {
	case <-ctx.Done():
		s.log.Error("request timed out")
		w.WriteHeader(http.StatusInternalServerError)
		w.Header().Add("Content-Type", "application/json")
		json.NewEncoder(w).Encode(serverResponse{
			Error: "request timed out",
		})
		return
	case resp := <-respch:
		if resp.err != nil {
			s.log.Error(resp.err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			w.Header().Add("Content-Type", "application/json")
			json.NewEncoder(w).Encode(serverResponse{
				Error: resp.err.Error(),
			})
			return
		}

		s.log.Info("Successfully created todo")

		w.Header().Add("Content-Type", "application/json")
		json.NewEncoder(w).Encode(serverResponse{
			Message: "Successfully created todo",
			Value:   &resp.val,
		})
		return
	}
}

func (*HttpServer) handleGetTodo(w http.ResponseWriter, r *http.Request) {
	panic("not implemented")
}

func (*HttpServer) handleUpdateTodo(w http.ResponseWriter, r *http.Request) {
	panic("not implemented")
}

func (*HttpServer) handleDeleteTodo(w http.ResponseWriter, r *http.Request) {
	panic("not implemented")
}
