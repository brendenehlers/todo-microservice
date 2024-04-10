package http

import (
	"context"
	"fmt"
	"net/http"

	"github.com/brendenehlers/todo-microservice"
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
			Handler: handler,
		},
		repo: config.Repo,
		ctx:  config.Ctx,
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
	ctx  context.Context
	log  todo.Logger
}

func (*HttpServer) Run() {
	panic("not implemented")
}

func (*HttpServer) Stop() {
	panic("not implemented")
}

func (*HttpServer) handleCreateTodo(w http.ResponseWriter, r *http.Request) {
	panic("not implemented")
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
