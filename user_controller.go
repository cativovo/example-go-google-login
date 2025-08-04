package main

import "context"

type getUserOutput struct {
	Body struct {
		User User `json:"user" doc:"Currently signed in user"`
	}
}

func getUser(store *SessionStore) humaHandler[struct{}, getUserOutput] {
	return func(ctx context.Context, input *struct{}) (*getUserOutput, error) {
		var resp getUserOutput
		resp.Body.User = userFromContext(ctx)
		return &resp, nil
	}
}
