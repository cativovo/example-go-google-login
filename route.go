package main

import (
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humachi"
	"github.com/go-chi/chi/v5"
)

func mountRoutes(r *chi.Mux, sessionStore *SessionStore, todoStore *TodoStore) {
	// TODO: remove me
	r.Group(func(r chi.Router) {
		r.Use(func(h http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				u, err := sessionStore.GetUser(r)
				if err != nil {
					http.Redirect(w, r, "/auth/login", http.StatusUnauthorized)
					return
				}

				r = r.WithContext(contextWithUser(r.Context(), u))
				h.ServeHTTP(w, r)
			})
		})
		r.Get("/", userPage())
	})

	mountAuthRoutes(r, sessionStore)
	mountDocsRoute(r)

	r.Group(func(r chi.Router) {
		// NOTE: to test the APIs, log in first at /auth/login to set the auth cookie
		// Scalar cannot set cookies - https://github.com/scalar/scalar/issues/3701
		api := huma.NewGroup(humachi.New(r, humaConfig), "/api")
		api.UseMiddleware(auth(api, sessionStore))

		mountTodoRoutes(api, todoStore)
		mountUserRoutes(api, sessionStore)
	})
}
