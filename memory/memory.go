package memory

import (
	"fmt"
	"time"

	"github.com/brendenehlers/todo-microservice/domain"
	"github.com/google/uuid"
)

var (
	ErrTodoDoesNotExist  = fmt.Errorf("todo does not exist")
	ErrTodoAlreadyExists = fmt.Errorf("todo already exists")
	ErrInvalidParameter  = fmt.Errorf("invalid function parameters")
)

func New(log domain.Logger) *InMemoryTodoRepository {
	return &InMemoryTodoRepository{
		todos: make(map[string]*domain.Todo),
		log:   log,
	}
}

type InMemoryTodoRepository struct {
	todos map[string]*domain.Todo
	log   domain.Logger
}

func (r *InMemoryTodoRepository) CreateTodo(newTodo *domain.NewTodo) (*domain.Todo, error) {
	if newTodo == nil {
		return nil, ErrInvalidParameter
	}

	// random uuidV4
	id := uuid.New().String()

	if _, ok := r.todos[id]; ok {
		return nil, ErrTodoAlreadyExists
	}

	todo := &domain.Todo{
		Id:          id,
		Description: newTodo.Description,
		CreatedAt:   time.Now(),
	}

	r.todos[id] = todo

	return todo, nil
}

func (r *InMemoryTodoRepository) GetTodo(id string) (*domain.Todo, error) {
	todo, ok := r.todos[id]
	if !ok {
		return nil, ErrTodoDoesNotExist
	}

	return todo, nil
}

func (r *InMemoryTodoRepository) GetTodos() (*[]domain.Todo, error) {
	todos := make([]domain.Todo, 0)

	for _, v := range r.todos {
		todos = append(todos, *v)
	}

	return &todos, nil
}

func (r *InMemoryTodoRepository) UpdateTodo(id string, todo *domain.UpdateTodo) (*domain.Todo, error) {
	if todo == nil {
		return nil, ErrInvalidParameter
	}

	if _, ok := r.todos[id]; !ok {
		return nil, ErrTodoDoesNotExist
	}

	r.todos[id].Done = todo.Done
	if todo.Done {
		r.todos[id].DoneAt = time.Now()
	}
	r.todos[id].Description = todo.Description
	r.todos[id].UpdatedAt = time.Now()

	return r.todos[id], nil
}

func (r *InMemoryTodoRepository) DeleteTodo(id string) error {
	delete(r.todos, id)

	return nil
}
