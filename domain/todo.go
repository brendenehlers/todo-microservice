package domain

import "time"

type Todo struct {
	Id          string    `json:"id"`
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
	GetTodo(id string) (*Todo, error)
	UpdateTodo(id string, todo *Todo) (*Todo, error)
	DeleteTodo(id string) error
}
