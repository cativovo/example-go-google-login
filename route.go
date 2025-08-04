package main

import (
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humachi"
	"github.com/go-chi/chi/v5"
	"github.com/markbates/goth/gothic"
)

func mountRoutes(r *chi.Mux, sessionStore *SessionStore, todoStore *TodoStore) {
	r.Route("/auth", func(r chi.Router) {
		r.Get("/begin", gothic.BeginAuthHandler)
		r.Get("/callback", authCallback(sessionStore))
		r.Get("/login", login())
		r.Get("/logout", logout(sessionStore))
	})

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
		r.Get("/", getUser())
	})

	r.Get("/docs", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(`<!doctype html>
<html>
  <head>
    <title>API Reference</title>
    <meta charset="utf-8" />
    <meta
      name="viewport"
      content="width=device-width, initial-scale=1" />
  </head>
  <body>
    <script
      id="api-reference"
      data-url="/openapi.yaml"></script>
    <script src="https://cdn.jsdelivr.net/npm/@scalar/api-reference"></script>
  </body>
</html>`))
	})

	r.Group(func(r chi.Router) {
		config := huma.DefaultConfig("My API", "1.0.0")
		config.Components.SecuritySchemes = map[string]*huma.SecurityScheme{
			"SessionAuth": {
				Type: "apiKey",
				In:   "cookie",
				Name: AuthSession,
			},
		}
		config.DocsPath = ""
		config.CreateHooks = nil
		security := []map[string][]string{
			{"SessionAuth": {}},
		}

		api := huma.NewGroup(humachi.New(r, config), "/api")
		api.UseMiddleware(auth(api, sessionStore))

		{
			todo := huma.NewGroup(api, "/todo")
			huma.Register(
				todo,
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
				todo,
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
				todo,
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
				todo,
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
	})
}
