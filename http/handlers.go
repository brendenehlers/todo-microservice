package http

import (
	"context"
	"encoding/json"
	"io"
	"net/http"

	"github.com/brendenehlers/todo-microservice/domain"
	"github.com/brendenehlers/todo-microservice/http/generated"
)

type GeneratedTodoRepository interface {
	CreateTodo(newTodo *generated.CreateTodoJSONRequestBody) (*generated.Todo, error)
	GetTodo(id *generated.TodoID) (*generated.Todo, error)
	GetTodos() (*[]generated.Todo, error)
	UpdateTodo(id *generated.TodoID, update *generated.UpdateTodoJSONRequestBody) (*generated.Todo, error)
	DeleteTodo(id *generated.TodoID) error
}

func newAPI(
	repo GeneratedTodoRepository,
	log domain.Logger,
) *api {
	return &api{
		repo: repo,
		log:  log,
	}
}

type api struct {
	repo GeneratedTodoRepository
	log  domain.Logger
}

func (a *api) handler() http.Handler {
	return generated.Handler(a)
}

type response struct {
	val any
	err error
}

func (*api) GetStatus(w http.ResponseWriter, r *http.Request) {
	status := "ok"
	json.NewEncoder(w).Encode(generated.Status{
		Status: &status,
	})
}

func (api *api) CreateTodo(w http.ResponseWriter, r *http.Request) {
	var newTodo generated.CreateTodoJSONRequestBody
	err := decodeRequestBody(r.Body, &newTodo)
	defer r.Body.Close()
	if err != nil {
		api.requestError(w, err)
		return
	}

	ctx, cancel, respch := processWithTimeout(r.Context(), func(respch chan response) {
		todo, err := api.repo.CreateTodo(&newTodo)
		respch <- response{
			val: todo,
			err: err,
		}
	})
	defer cancel()

	select {
	case <-ctx.Done():
		api.requestTimeout(w)
		return
	case resp := <-respch:
		if resp.err != nil {
			api.requestError(w, resp.err)
			return
		}

		api.log.Info("Successfully created todo")
		api.sendTodoResponse(w, resp.val.(*generated.Todo))
		return
	}
}

func (api *api) GetTodo(w http.ResponseWriter, r *http.Request, todoId generated.TodoID) {
	ctx, cancel, respch := processWithTimeout(r.Context(), func(respch chan response) {
		val, err := api.repo.GetTodo(&todoId)
		respch <- response{
			val: val,
			err: err,
		}
	})
	defer cancel()

	select {
	case <-ctx.Done():
		api.requestTimeout(w)
		return
	case resp := <-respch:
		if resp.err != nil {
			api.requestError(w, resp.err)
			return
		}

		api.log.Info("Successfully found todo")
		api.sendTodoResponse(w, resp.val.(*generated.Todo))
		return
	}
}

func (api *api) GetTodos(w http.ResponseWriter, r *http.Request) {
	ctx, cancel, respch := processWithTimeout(r.Context(), func(respch chan response) {
		val, err := api.repo.GetTodos()
		respch <- response{
			val: val,
			err: err,
		}
	})
	defer cancel()

	select {
	case <-ctx.Done():
		api.requestTimeout(w)
		return
	case resp := <-respch:
		if resp.err != nil {
			api.requestError(w, resp.err)
			return
		}

		api.log.Info("Successfully retrieved todos")
		api.sendTodosResponse(w, resp.val.(*[]generated.Todo))
	}
}

func (api *api) UpdateTodo(w http.ResponseWriter, r *http.Request, todoId generated.TodoID) {
	var todo generated.UpdateTodoJSONRequestBody
	err := decodeRequestBody(r.Body, &todo)
	if err != nil {
		api.requestError(w, err)
		return
	}

	ctx, cancel, respch := processWithTimeout(r.Context(), func(respch chan response) {
		val, err := api.repo.UpdateTodo(&todoId, &todo)
		respch <- response{
			val: val,
			err: err,
		}
	})
	defer cancel()

	select {
	case <-ctx.Done():
		api.requestTimeout(w)
		return
	case resp := <-respch:
		if resp.err != nil {
			api.requestError(w, resp.err)
			return
		}

		api.log.Info("Successfully updated todo")
		api.sendTodoResponse(w, resp.val.(*generated.Todo))
	}
}

func (api *api) DeleteTodo(w http.ResponseWriter, r *http.Request, todoId generated.TodoID) {
	ctx, cancel, respch := processWithTimeout(r.Context(), func(respch chan response) {
		err := api.repo.DeleteTodo(&todoId)
		respch <- response{
			err: err,
		}
	})
	defer cancel()

	select {
	case <-ctx.Done():
		api.requestTimeout(w)
		return
	case resp := <-respch:
		if resp.err != nil {
			api.requestError(w, resp.err)
		}

		msg := "Successfully deleted todo"
		api.log.Info(msg)
		api.requestSuccessWithMessage(w, &msg)
	}
}

func (api *api) requestTimeout(w http.ResponseWriter) {
	errStr := ErrRequestTimedOut.Error()
	api.log.Error(ErrRequestTimedOut.Error())

	w.WriteHeader(http.StatusRequestTimeout)
	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(generated.Error{
		Error: &errStr,
	})
}

func (api *api) requestError(w http.ResponseWriter, err error) {
	errStr := err.Error()
	api.log.Error(errStr)

	w.WriteHeader(http.StatusInternalServerError)
	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(generated.Error{
		Error: &errStr,
	})
}

func (api *api) sendTodoResponse(w http.ResponseWriter, todo *generated.Todo) {
	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(generated.TodoResponse{
		Value: todo,
	})
}

func (api *api) sendTodosResponse(w http.ResponseWriter, todos *[]generated.Todo) {
	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(generated.TodosResponse{
		Value: todos,
	})
}

func (api *api) requestSuccessWithMessage(w http.ResponseWriter, message *string) {
	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(generated.MessageResponse{
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
