package main

import (
	"github.com/go-chi/chi/v5"
	"github.com/markbates/goth/gothic"
)

func mountAuthRoutes(r *chi.Mux, sessionStore *SessionStore) {
	r.Route("/auth", func(r chi.Router) {
		r.Get("/begin", gothic.BeginAuthHandler)
		r.Get("/callback", authCallback(sessionStore))
		r.Get("/login", login())
		r.Get("/logout", logout(sessionStore))
	})
}
