package memory

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/brendenehlers/todo-microservice/domain"
)

var (
	ErrTodoDoesNotExist  = fmt.Errorf("todo does not exist")
	ErrTodoAlreadyExists = fmt.Errorf("todo already exists")
	ErrInvalidParameter  = fmt.Errorf("invalid function parameters")
)

func New(log domain.Logger) *InMemoryTodoRepository {
	return &InMemoryTodoRepository{
		todos: make(map[int]*domain.Todo),
		log:   log,
	}
}

type InMemoryTodoRepository struct {
	todos map[int]*domain.Todo
	log   domain.Logger
}

func (r *InMemoryTodoRepository) CreateTodo(newTodo *domain.NewTodo) (*domain.Todo, error) {
	if newTodo == nil {
		return nil, ErrInvalidParameter
	}

	// TODO make this better
	id := rand.Intn(10000)

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

func (r *InMemoryTodoRepository) GetTodo(id int) (*domain.Todo, error) {
	todo, ok := r.todos[id]
	if !ok {
		return nil, ErrTodoDoesNotExist
	}

	return todo, nil
}

func (r *InMemoryTodoRepository) UpdateTodo(id int, todo *domain.Todo) (*domain.Todo, error) {
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

func (r *InMemoryTodoRepository) DeleteTodo(id int) error {
	delete(r.todos, id)

	return nil
}
