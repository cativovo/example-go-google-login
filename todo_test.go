package main

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTodoStore_CreateTodo(t *testing.T) {
	user := User{
		ID:        "123",
		FirstName: "Juan",
		LastName:  "Usa",
		Email:     "juanusa@example.com",
		AvatarURL: "https://www.images.com/123",
	}
	ctxWithUser := contextWithUser(context.Background(), user)
	tests := []struct {
		name  string
		input CreateTodoInput
		want  Todo
	}{
		{
			name: "create todo",
			input: CreateTodoInput{
				Name:        "todo 1 name",
				Description: "todo 1 description",
			},
			want: Todo{
				UserID:      user.ID,
				Name:        "todo 1 name",
				Description: "todo 1 description",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ts := NewTodoStore()
			created := ts.CreateTodo(ctxWithUser, tt.input)
			assert.NotEmpty(t, created.ID)
			assert.NotZero(t, created.CreatedAt)
			assert.NotZero(t, created.UpdatedAt)
			assert.Equal(t, tt.want.Name, created.Name)
			assert.Equal(t, tt.want.Description, created.Description)

			got, err := ts.GetTodo(ctxWithUser, created.ID)
			assert.NoError(t, err)
			assert.Equal(t, created, got)
		})
	}

	t.Run("test access", func(t *testing.T) {
		input := CreateTodoInput{
			Name:        "todo 1 name",
			Description: "todo 1 description",
		}
		ts := NewTodoStore()
		created := ts.CreateTodo(ctxWithUser, input)

		ctxWithOtherUser := contextWithUser(context.Background(), User{
			ID:        "2",
			FirstName: "other",
			LastName:  "user",
			Email:     "otheruser@gmail.com",
			AvatarURL: "https://www.images/123",
		})
		_, err := ts.GetTodo(ctxWithOtherUser, created.ID)
		assert.Error(t, err)
	})
}

func TestTodoStore_ListTodos(t *testing.T) {
	user := User{
		ID:        "123",
		FirstName: "Juan",
		LastName:  "Usa",
		Email:     "juanusa@example.com",
		AvatarURL: "https://www.images.com/123",
	}
	ctxWithUser := contextWithUser(context.Background(), user)
	t.Run("get todos", func(t *testing.T) {
		ts := NewTodoStore()
		created1 := ts.CreateTodo(ctxWithUser, CreateTodoInput{
			Name:        "todo name 1",
			Description: "todo description 2",
		})
		created2 := ts.CreateTodo(ctxWithUser, CreateTodoInput{
			Name:        "todo name 2",
			Description: "todo description 2",
		})
		got := ts.ListTodos(ctxWithUser)
		assert.Len(t, got, 2)
		assert.Contains(t, got, created1)
		assert.Contains(t, got, created2)
	})

	t.Run("invalid access", func(t *testing.T) {
		ts := NewTodoStore()
		ts.CreateTodo(ctxWithUser, CreateTodoInput{
			Name:        "todo name 1",
			Description: "todo description 2",
		})
		ts.CreateTodo(ctxWithUser, CreateTodoInput{
			Name:        "todo name 2",
			Description: "todo description 2",
		})

		ctxWithOtherUser := contextWithUser(context.Background(), User{
			ID:        "2",
			FirstName: "other",
			LastName:  "user",
			Email:     "otheruser@gmail.com",
			AvatarURL: "https://www.images/123",
		})
		got := ts.ListTodos(ctxWithOtherUser)
		assert.Len(t, got, 0)
	})
}

func TestTodoStore_UpdateTodo(t *testing.T) {
	user := User{
		ID:        "123",
		FirstName: "Juan",
		LastName:  "Usa",
		Email:     "juanusa@example.com",
		AvatarURL: "https://www.images.com/123",
	}
	ctxWithUser := contextWithUser(context.Background(), user)
	tests := []struct {
		name       string
		todoCreate CreateTodoInput
		todoUpdate UpdateTodoInput
	}{
		{
			name: "update todo",
			todoCreate: CreateTodoInput{
				Name:        "todo name 1",
				Description: "todo description 1",
			},
			todoUpdate: UpdateTodoInput{
				Name:        toPtr("todo name updated 1"),
				Description: toPtr("todo description updated 1"),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ts := NewTodoStore()
			created := ts.CreateTodo(ctxWithUser, tt.todoCreate)
			tt.todoUpdate.ID = created.ID
			updated, err := ts.UpdateTodo(ctxWithUser, tt.todoUpdate)
			assert.NoError(t, err)
			assert.NotZero(t, updated.CreatedAt)
			assert.NotZero(t, updated.UpdatedAt)
			assert.NotEqual(t, created, updated)
			assert.Equal(t, *tt.todoUpdate.Name, updated.Name)
			assert.Equal(t, *tt.todoUpdate.Description, updated.Description)
		})
	}

	t.Run("invalid access", func(t *testing.T) {
		ts := NewTodoStore()
		created := ts.CreateTodo(ctxWithUser, CreateTodoInput{
			Name:        "todo name 1",
			Description: "todo description 2",
		})

		ctxWithOtherUser := contextWithUser(context.Background(), User{
			ID:        "2",
			FirstName: "other",
			LastName:  "user",
			Email:     "otheruser@gmail.com",
			AvatarURL: "https://www.images/123",
		})
		_, err := ts.UpdateTodo(ctxWithOtherUser, UpdateTodoInput{
			ID:   created.ID,
			Name: toPtr("updated name"),
		})
		assert.Error(t, err)
	})
}

func TestTodoStore_DeleteTodo(t *testing.T) {
	user := User{
		ID:        "123",
		FirstName: "Juan",
		LastName:  "Usa",
		Email:     "juanusa@example.com",
		AvatarURL: "https://www.images.com/123",
	}
	ctxWithUser := contextWithUser(context.Background(), user)
	tests := []struct {
		name       string
		todoCreate CreateTodoInput
	}{
		{
			name: "delete todo",
			todoCreate: CreateTodoInput{
				Name:        "todo name 1",
				Description: "todo description 1",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ts := NewTodoStore()
			created := ts.CreateTodo(ctxWithUser, tt.todoCreate)
			err := ts.DeleteTodo(ctxWithUser, created.ID)
			assert.NoError(t, err)

			_, err = ts.GetTodo(ctxWithUser, created.ID)
			assert.Error(t, err)
		})
	}

	t.Run("invalid access", func(t *testing.T) {
		ts := NewTodoStore()
		created := ts.CreateTodo(ctxWithUser, CreateTodoInput{
			Name:        "todo name 1",
			Description: "todo description 2",
		})

		ctxWithOtherUser := contextWithUser(context.Background(), User{
			ID:        "2",
			FirstName: "other",
			LastName:  "user",
			Email:     "otheruser@gmail.com",
			AvatarURL: "https://www.images/123",
		})
		err := ts.DeleteTodo(ctxWithOtherUser, created.ID)
		assert.Error(t, err)
	})
}

func toPtr[T any](v T) *T {
	return &v
}
