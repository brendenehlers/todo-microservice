// Package generated provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/deepmap/oapi-codegen/v2 version v2.1.0 DO NOT EDIT.
package generated

import (
	"time"

	openapi_types "github.com/oapi-codegen/runtime/types"
)

// Error defines model for Error.
type Error struct {
	Error *string `json:"error,omitempty"`
}

// MessageResponse defines model for MessageResponse.
type MessageResponse struct {
	Message *string `json:"message,omitempty"`
}

// Status defines model for Status.
type Status struct {
	Status *string `json:"status,omitempty"`
}

// Todo defines model for Todo.
type Todo struct {
	CreatedAt   *time.Time          `json:"createdAt,omitempty"`
	Description *string             `json:"description,omitempty"`
	Done        *bool               `json:"done,omitempty"`
	DoneAt      *time.Time          `json:"doneAt,omitempty"`
	Id          *openapi_types.UUID `json:"id,omitempty"`
	UpdatedAt   *time.Time          `json:"updatedAt,omitempty"`
}

// TodoResponse defines model for TodoResponse.
type TodoResponse struct {
	Message *string `json:"message,omitempty"`
	Value   *Todo   `json:"value,omitempty"`
}

// TodoID defines model for TodoID.
type TodoID = openapi_types.UUID

// N500 defines model for 500.
type N500 = Error

// CreateTodoJSONBody defines parameters for CreateTodo.
type CreateTodoJSONBody struct {
	Description *string `json:"description,omitempty"`
}

// CreateTodoJSONRequestBody defines body for CreateTodo for application/json ContentType.
type CreateTodoJSONRequestBody CreateTodoJSONBody