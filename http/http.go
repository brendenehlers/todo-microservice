package http

import (
	"net/http"

	"github.com/brendenehlers/todo-microservice"
)

func CreateHTTPServer(todos todo.TodoRepository) *HttpServer {
	return &HttpServer{
		todos: todos,
	}
}

type HttpServer struct {
	todos todo.TodoRepository
}

func (*HttpServer) HandleCreateTodo(w http.ResponseWriter, r *http.Request) {
	panic("not implemented")
}

func (*HttpServer) HandleGetTodo(w http.ResponseWriter, r *http.Request) {
	panic("not implemented")
}

func (*HttpServer) HandleUpdateTodo(w http.ResponseWriter, r *http.Request) {
	panic("not implemented")
}

func (*HttpServer) HandleDeleteTodo(w http.ResponseWriter, r *http.Request) {
	panic("not implemented")
}
