package main

import (
	"context"
	"net/http"
)

type contextKey string

const contextKeyUser contextKey = "user"

func contextWithUser(ctx context.Context, u User) context.Context {
	return context.WithValue(ctx, contextKeyUser, u)
}

func userFromContext(ctx context.Context) User {
	u, ok := ctx.Value(contextKeyUser).(User)
	if !ok {
		panic("User not found")
	}
	return u
}

func auth(store *SessionStore) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			u, err := store.GetUser(r)
			if err != nil {
				http.Redirect(w, r, "/auth/login", http.StatusFound)
				return
			}

			r = r.WithContext(contextWithUser(r.Context(), u))
			next.ServeHTTP(w, r)
		})
	}
}
