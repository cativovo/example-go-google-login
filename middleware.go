package main

import (
	"context"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humachi"
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

type humaMiddleware func(ctx huma.Context, next func(huma.Context))

func auth(api huma.API, store *SessionStore) humaMiddleware {
	return func(ctx huma.Context, next func(huma.Context)) {
		r, _ := humachi.Unwrap(ctx)
		u, err := store.GetUser(r)
		if err != nil {
			huma.WriteErr(api, ctx, http.StatusUnauthorized, "Unauthorized")
			return
		}

		ctx = huma.WithContext(ctx, contextWithUser(r.Context(), u))
		next(ctx)
	}
}
