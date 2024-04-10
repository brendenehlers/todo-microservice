package todo

import "time"

type Todo struct {
	Id          int       `json:"id"`
	Done        bool      `json:"done"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
	DoneAt      time.Time `json:"doneAt"`
}

type NewTodo struct {
	Description string `json:"description"`
}

type TodoRepository interface {
	CreateTodo(newTodo *NewTodo) (*Todo, error)
	GetTodo(id int) (*Todo, error)
	UpdateTodo(todo *Todo) (*Todo, error)
	DeleteTodo(id int) error
}
