package memory

import "github.com/brendenehlers/todo-microservice"

func New() *InMemoryTodoRepository {
	return &InMemoryTodoRepository{
		todos: make(map[string]*todo.Todo),
	}
}

type InMemoryTodoRepository struct {
	todos map[string]*todo.Todo
}

func (t *InMemoryTodoRepository) CreateTodo(newTodo *todo.NewTodo) (*todo.Todo, error) {
	panic("not implemented")
}

func (t *InMemoryTodoRepository) GetTodo(id string) (*todo.Todo, error) {
	panic("not implemented")
}

func (t *InMemoryTodoRepository) UpdateTodo(todo *todo.Todo) (*todo.Todo, error) {
	panic("not implemented")
}

func (t *InMemoryTodoRepository) DeleteTodo(id string) error {
	panic("not implemented")
}
