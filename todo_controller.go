package main

import (
	"context"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
)

type handler[I, O any] func(ctx context.Context, input *I) (*O, error)

type getTodoInput struct {
	ID string `path:"id" doc:"ID of the todo to get" example:"123"`
}
type getTodoOutput struct {
	Body struct {
		Todo Todo `json:"todo" doc:"A todo"`
	}
}

func getTodo(store *TodoStore) handler[getTodoInput, getTodoOutput] {
	return func(ctx context.Context, input *getTodoInput) (*getTodoOutput, error) {
		todo, err := store.GetTodo(ctx, input.ID)
		if err != nil {
			return nil, huma.NewError(http.StatusNotFound, err.Error())
		}
		resp := new(getTodoOutput)
		resp.Body.Todo = todo
		return resp, nil
	}
}

type createTodoInput struct {
	Body struct {
		Name        string `json:"name" doc:"Name of the todo" example:"name"`
		Description string `json:"description" doc:"Description of todo" example:"description"`
	}
}
type createTodoOutput struct {
	Body struct {
		Todo Todo `json:"todo" doc:"The created todo"`
	}
}

func createTodo(store *TodoStore) handler[createTodoInput, createTodoOutput] {
	return func(ctx context.Context, input *createTodoInput) (*createTodoOutput, error) {
		todo := store.CreateTodo(ctx, TodoCreate(input.Body))
		resp := new(createTodoOutput)
		resp.Body.Todo = todo
		return resp, nil
	}
}

type updateTodoInput struct {
	ID   string `path:"id" doc:"ID of the todo to update" example:"123"`
	Body struct {
		Name        *string `json:"name,omitempty" doc:"Name of the todo" example:"updated name"`
		Description *string `json:"description,omitempty" doc:"Description of the todo" example:"updated description"`
	}
}
type updateTodoOutput struct {
	Body struct {
		Todo Todo `json:"todo" doc:"The updated todo"`
	}
}

func updateTodo(store *TodoStore) handler[updateTodoInput, updateTodoOutput] {
	return func(ctx context.Context, input *updateTodoInput) (*updateTodoOutput, error) {
		todo, err := store.UpdateTodo(ctx, TodoUpdate{
			ID:          input.ID,
			Name:        input.Body.Name,
			Description: input.Body.Description,
		})
		if err != nil {
			return nil, huma.NewError(http.StatusNotFound, err.Error())
		}
		resp := new(updateTodoOutput)
		resp.Body.Todo = todo
		return resp, nil
	}
}

type deleteTodoInput struct {
	ID string `path:"id" doc:"ID of the todo to delete" example:"123"`
}

func deleteTodo(store *TodoStore) handler[deleteTodoInput, struct{}] {
	return func(ctx context.Context, input *deleteTodoInput) (*struct{}, error) {
		err := store.DeleteTodo(ctx, input.ID)
		if err != nil {
			return nil, huma.NewError(http.StatusNotFound, err.Error())
		}
		return nil, nil
	}
}
