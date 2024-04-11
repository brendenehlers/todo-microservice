package http

import (
	"github.com/brendenehlers/todo-microservice/domain"
	"github.com/brendenehlers/todo-microservice/http/generated"
	"github.com/google/uuid"
)

func newAdapter(repo domain.TodoRepository) *adapter {
	return &adapter{
		repo: repo,
	}
}

type adapter struct {
	repo domain.TodoRepository
}

func (a *adapter) CreateTodo(newTodo *generated.CreateTodoJSONRequestBody) (*generated.Todo, error) {
	domainNewTodo := convertGeneratedNewTodoToDomainNewTodo(newTodo)

	domainTodo, err := a.repo.CreateTodo(domainNewTodo)
	if err != nil {
		return nil, err
	}

	return covertDomainTodoToGeneratedTodo(domainTodo)
}

func (a *adapter) GetTodo(id *generated.TodoID) (*generated.Todo, error) {
	idStr := id.String()

	domainTodo, err := a.repo.GetTodo(idStr)
	if err != nil {
		return nil, err
	}

	return covertDomainTodoToGeneratedTodo(domainTodo)
}

func (a *adapter) GetTodos() (*[]generated.Todo, error) {
	domainTodos, err := a.repo.GetTodos()
	if err != nil {
		return nil, err
	}

	todos := make([]generated.Todo, 0)
	for _, dTodo := range *domainTodos {
		todo, err := covertDomainTodoToGeneratedTodo(&dTodo)
		if err != nil {
			return nil, err
		}
		todos = append(todos, *todo)
	}

	return &todos, nil
}

func (a *adapter) UpdateTodo(id *generated.TodoID, update *generated.UpdateTodoJSONRequestBody) (*generated.Todo, error) {
	idStr := id.String()
	domainUpdateTodo := convertGeneratedUpdateTodoToDomainUpdateTodo(update)

	todo, err := a.repo.UpdateTodo(idStr, domainUpdateTodo)
	if err != nil {
		return nil, err
	}

	return covertDomainTodoToGeneratedTodo(todo)
}

func (a *adapter) DeleteTodo(id *generated.TodoID) error {
	idStr := id.String()

	return a.repo.DeleteTodo(idStr)
}

func convertGeneratedNewTodoToDomainNewTodo(newTodo *generated.CreateTodoJSONRequestBody) *domain.NewTodo {
	return &domain.NewTodo{
		Description: *newTodo.Description,
	}
}

func covertDomainTodoToGeneratedTodo(todo *domain.Todo) (*generated.Todo, error) {
	uuidObj, err := uuid.Parse(todo.Id)
	if err != nil {
		return nil, err
	}

	return &generated.Todo{
		Id:          &uuidObj,
		Done:        &todo.Done,
		Description: &todo.Description,
		DoneAt:      &todo.DoneAt,
		CreatedAt:   &todo.CreatedAt,
		UpdatedAt:   &todo.UpdatedAt,
	}, nil
}

func convertGeneratedUpdateTodoToDomainUpdateTodo(todo *generated.UpdateTodoJSONRequestBody) *domain.UpdateTodo {
	return &domain.UpdateTodo{
		Done:        *todo.Done,
		Description: *todo.Description,
	}
}
