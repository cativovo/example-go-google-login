package main

import (
	"context"
	"errors"
	"strconv"
	"sync"
	"time"
)

type Todo struct {
	ID          string    `json:"id" doc:"Unique identifier" example:"123"`
	Name        string    `json:"name" doc:"Name" example:"Read"`
	Description string    `json:"description" doc:"Description" example:"Read the book of love"`
	CreatedAt   time.Time `json:"created_at" doc:"Created date"`
	UpdatedAt   time.Time `json:"updated_at" doc:"Last modified date"`
	UserID      string    `json:"-"`
}

type TodoStore struct {
	muStore sync.RWMutex
	store   map[string]Todo

	muID sync.Mutex
	id   int
}

func NewTodoStore() *TodoStore {
	return &TodoStore{
		store: make(map[string]Todo),
	}
}

func (ts *TodoStore) getID() string {
	ts.muID.Lock()
	defer ts.muID.Unlock()
	ts.id++
	return strconv.Itoa(ts.id)
}

type TodoCreate struct {
	Name        string
	Description string
}

func (ts *TodoStore) CreateTodo(ctx context.Context, input TodoCreate) Todo {
	u := userFromContext(ctx)
	now := time.Now()
	t := Todo{
		ID:          ts.getID(),
		UserID:      u.ID,
		Name:        input.Name,
		Description: input.Description,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
	ts.muStore.Lock()
	defer ts.muStore.Unlock()
	ts.store[t.ID] = t
	return t
}

func (ts *TodoStore) GetTodo(ctx context.Context, id string) (Todo, error) {
	ts.muStore.RLock()
	defer ts.muStore.RUnlock()

	t, ok := ts.store[id]
	if !ok {
		return Todo{}, errors.New("todo not found")
	}
	u := userFromContext(ctx)
	if t.UserID != u.ID {
		return Todo{}, errors.New("todo not found")
	}
	return t, nil
}

func (ts *TodoStore) ListTodos(ctx context.Context) []Todo {
	ts.muStore.RLock()
	defer ts.muStore.RUnlock()

	todos := make([]Todo, 0)
	u := userFromContext(ctx)
	for _, v := range ts.store {
		if v.UserID == u.ID {
			todos = append(todos, v)
		}
	}
	return todos
}

type TodoUpdate struct {
	ID          string
	Name        *string
	Description *string
}

func (ts *TodoStore) UpdateTodo(ctx context.Context, input TodoUpdate) (Todo, error) {
	t, err := ts.GetTodo(ctx, input.ID)
	if err != nil {
		return Todo{}, err
	}

	if input.Name != nil {
		t.Name = *input.Name
	}

	if input.Description != nil {
		t.Description = *input.Description
	}

	ts.muStore.Lock()
	defer ts.muStore.Unlock()
	ts.store[t.ID] = t
	return t, nil
}

func (ts *TodoStore) DeleteTodo(ctx context.Context, id string) error {
	_, err := ts.GetTodo(ctx, id)
	if err != nil {
		return err
	}

	ts.muStore.Lock()
	defer ts.muStore.Unlock()
	delete(ts.store, id)
	return nil
}
