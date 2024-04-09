package todo

type Todo struct {
	Id          string `json:"id"`
	Done        bool   `json:"done"`
	Description string `json:"description"`
	CreatedAt   string `json:"createdAt"`
	UpdatedAt   string `json:"updatedAt"`
	DoneAt      string `json:"doneAt"`
}

type NewTodo struct {
	Done        bool   `json:"done"`
	Description string `json:"description"`
}

type TodoRepository interface {
	CreateTodo(newTodo *NewTodo) (*Todo, error)
	GetTodo(id string) (*Todo, error)
	UpdateTodo(todo *Todo) (*Todo, error)
	DeleteTodo(id string) error
}
