package memory

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/brendenehlers/todo-microservice/domain"
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
	// TODO make this better
	id := rand.Intn(10000)

	if _, ok := r.todos[id]; ok {
		return nil, fmt.Errorf("unexpected error with creating todo")
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
		return nil, fmt.Errorf("unable to find todo")
	}

	return todo, nil
}

func (r *InMemoryTodoRepository) UpdateTodo(id int, todo *domain.Todo) (*domain.Todo, error) {
	if _, ok := r.todos[id]; !ok {
		return nil, fmt.Errorf("todo with given id does not exist")
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
