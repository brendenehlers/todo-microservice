package http

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"strconv"
	"time"

	"github.com/brendenehlers/todo-microservice/domain"
)

const (
	REQUEST_TIMEOUT = time.Millisecond * 200
)

type HTTPServerConfig struct {
	Addr string
	Repo domain.TodoRepository
	Ctx  context.Context
	Log  domain.Logger
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
	repo domain.TodoRepository
	log  domain.Logger
}

type serverResponse struct {
	Message string      `json:"message,omitempty"`
	Value   interface{} `json:"value,omitempty"`
	Error   string      `json:"error,omitempty"`
}

type response struct {
	val *domain.Todo
	err error
}

func (s *HttpServer) Run() {
	s.log.Info(fmt.Sprintf("Server running on %s", s.Addr))
	s.ListenAndServe()
}

func (*HttpServer) Stop() {
	panic("not implemented")
}

func (s *HttpServer) handleCreateTodo(w http.ResponseWriter, r *http.Request) {
	var newTodo domain.NewTodo
	err := decodeRequestBody(r.Body, &newTodo)
	if err != nil {
		s.requestError(w, err)
		return
	}

	ctx, cancel, respch := processWithTimeout(r.Context(), func(respch chan response) {
		todo, err := s.repo.CreateTodo(&newTodo)
		respch <- response{
			val: todo,
			err: err,
		}
	})
	defer cancel()

	select {
	case <-ctx.Done():
		s.requestTimeout(w)
		return
	case resp := <-respch:
		if resp.err != nil {
			s.requestError(w, resp.err)
			return
		}

		s.log.Info("Successfully created todo")
		s.requestSuccess(w, resp)
		return
	}
}

func (s *HttpServer) handleGetTodo(w http.ResponseWriter, r *http.Request) {
	todoId, err := getTodoIdFromRequest(r)
	if err != nil {
		s.requestError(w, err)
		return
	}

	ctx, cancel, respch := processWithTimeout(r.Context(), func(respch chan response) {
		val, err := s.repo.GetTodo(todoId)
		respch <- response{
			val: val,
			err: err,
		}
	})
	defer cancel()

	select {
	case <-ctx.Done():
		s.requestTimeout(w)
		return
	case resp := <-respch:
		if resp.err != nil {
			s.requestError(w, resp.err)
			return
		}

		s.log.Info("Successfully found todo")
		s.requestSuccess(w, resp)
		return
	}
}

func (s *HttpServer) handleUpdateTodo(w http.ResponseWriter, r *http.Request) {
	todoId, err := getTodoIdFromRequest(r)
	if err != nil {
		s.requestError(w, err)
		return
	}

	var todo domain.Todo
	err = decodeRequestBody(r.Body, &todo)
	if err != nil {
		s.requestError(w, err)
		return
	}

	ctx, cancel, respch := processWithTimeout(r.Context(), func(respch chan response) {
		val, err := s.repo.UpdateTodo(todoId, &todo)
		respch <- response{
			val: val,
			err: err,
		}
	})
	defer cancel()

	select {
	case <-ctx.Done():
		s.requestTimeout(w)
		return
	case resp := <-respch:
		if resp.err != nil {
			s.requestError(w, resp.err)
			return
		}

		s.log.Info("Updated todo successfully")
		s.requestSuccess(w, resp)
	}
}

func (s *HttpServer) handleDeleteTodo(w http.ResponseWriter, r *http.Request) {
	todoId, err := getTodoIdFromRequest(r)
	if err != nil {
		s.requestError(w, err)
	}

	ctx, cancel, respch := processWithTimeout(r.Context(), func(respch chan response) {
		err := s.repo.DeleteTodo(todoId)
		respch <- response{
			err: err,
		}
	})
	defer cancel()

	select {
	case <-ctx.Done():
		s.requestTimeout(w)
		return
	case resp := <-respch:
		if resp.err != nil {
			s.requestError(w, resp.err)
		}

		s.log.Info("Successfully deleted todo")
		s.requestSuccessWithMessage(w, "Successfully deleted todo")
	}
}

func (s *HttpServer) requestTimeout(w http.ResponseWriter) {
	s.log.Error("request timed out")
	w.WriteHeader(http.StatusInternalServerError)
	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(serverResponse{
		Error: "request timed out",
	})
}

func (s *HttpServer) requestError(w http.ResponseWriter, err error) {
	s.log.Error(err.Error())
	w.WriteHeader(http.StatusInternalServerError)
	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(serverResponse{
		Error: err.Error(),
	})
}

func (s *HttpServer) requestSuccess(w http.ResponseWriter, resp response) {
	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(serverResponse{
		Value: resp.val,
	})
}

func (s *HttpServer) requestSuccessWithMessage(w http.ResponseWriter, message string) {
	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(serverResponse{
		Message: message,
	})
}

func decodeRequestBody(r io.ReadCloser, data any) error {
	err := json.NewDecoder(r).Decode(data)
	defer r.Close()
	return err
}

func processWithTimeout(parentCtx context.Context, fn func(respch chan response)) (context.Context, context.CancelFunc, chan response) {
	ctx, cancel := context.WithTimeout(parentCtx, REQUEST_TIMEOUT)
	respch := make(chan response)

	go fn(respch)

	return ctx, cancel, respch
}

func getTodoIdFromRequest(r *http.Request) (int, error) {
	return strconv.Atoi(r.PathValue("todoId"))
}
