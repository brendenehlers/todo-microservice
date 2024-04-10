package memory

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/brendenehlers/todo-microservice"
)

func New(log todo.Logger) *InMemoryTodoRepository {
	return &InMemoryTodoRepository{
		todos: make(map[int]*todo.Todo),
		log:   log,
	}
}

type InMemoryTodoRepository struct {
	todos map[int]*todo.Todo
	log   todo.Logger
}

func (r *InMemoryTodoRepository) CreateTodo(newTodo *todo.NewTodo) (*todo.Todo, error) {
	// TODO make this better
	id := rand.Intn(10000)

	if _, ok := r.todos[id]; ok {
		return nil, fmt.Errorf("unexpected error with creating todo")
	}

	todo := &todo.Todo{
		Id:          id,
		Description: newTodo.Description,
		CreatedAt:   time.Now(),
	}

	r.todos[id] = todo

	return todo, nil
}

func (r *InMemoryTodoRepository) GetTodo(id int) (*todo.Todo, error) {
	todo, ok := r.todos[id]
	if !ok {
		return nil, fmt.Errorf("unable to find todo")
	}

	return todo, nil
}

func (r *InMemoryTodoRepository) UpdateTodo(todo *todo.Todo) (*todo.Todo, error) {
	if _, ok := r.todos[todo.Id]; !ok {
		return nil, fmt.Errorf("todo with given id does not exist")
	}

	r.todos[todo.Id] = todo

	return todo, nil
}

func (r *InMemoryTodoRepository) DeleteTodo(id int) error {
	delete(r.todos, id)

	return nil
}
