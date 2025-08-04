package main

import (
	"net/http"

	"github.com/danielgtaylor/huma/v2"
)

func mountTodoRoutes(api huma.API, todoStore *TodoStore) {
	todoGroup := huma.NewGroup(api, "/todo")
	huma.Register(
		todoGroup,
		huma.Operation{
			OperationID: "get-todo",
			Method:      http.MethodGet,
			Path:        "/{id}",
			Summary:     "Get a todo by id",
			Security:    security,
		},
		getTodo(todoStore),
	)
	huma.Register(
		todoGroup,
		huma.Operation{
			OperationID:   "create-todo",
			Method:        http.MethodPost,
			Path:          "",
			Summary:       "Create a todo",
			DefaultStatus: http.StatusCreated,
			Security:      security,
		},
		createTodo(todoStore),
	)
	huma.Register(
		todoGroup,
		huma.Operation{
			OperationID: "update-todo",
			Method:      http.MethodPatch,
			Path:        "/{id}",
			Summary:     "Update a todo",
			Security:    security,
		},
		updateTodo(todoStore),
	)
	huma.Register(
		todoGroup,
		huma.Operation{
			OperationID: "delete-todo",
			Method:      http.MethodDelete,
			Path:        "/{id}",
			Summary:     "Delete a todo",
			Security:    security,
		},
		deleteTodo(todoStore),
	)
}
