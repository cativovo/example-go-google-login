package main

import (
	"context"
	"encoding/gob"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/google"
)

func setGothicConfig(store *SessionStore, cfg config) {
	scopes := []string{"profile", "email"}
	goth.UseProviders(
		google.New(os.Getenv("GOOGLE_KEY"), os.Getenv("GOOGLE_SECRET"), "http://127.0.0.1:8000/auth/callback", scopes...),
	)
	gothic.Store = store
	gothic.GetProviderName = func(_ *http.Request) (string, error) {
		return "google", nil
	}
}

func main() {
	// fix for -> securecookie: error - caused by: securecookie: error - caused by: gob: type not registered for interface: main.User
	gob.Register(User{})

	cfg, err := getConfig()
	if err != nil {
		log.Fatal("Failed to get the config: ", err)
	}

	sessionStore := NewSessionStore()
	setGothicConfig(sessionStore, cfg)

	todoStore := NewTodoStore()
	router := chi.NewRouter()
	mountRoutes(router, sessionStore, todoStore)

	server := http.Server{
		Addr:    ":8000",
		Handler: router,
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		ctx, cancel := signal.NotifyContext(ctx, os.Interrupt, os.Kill)
		defer cancel()

		<-ctx.Done()

		ctx, cancel = context.WithTimeout(ctx, time.Second*10)
		defer cancel()
		log.Println("shutting down the server")
		if err := server.Shutdown(ctx); err != nil {
			log.Fatal(err)
		}
	}()

	log.Println("listening on localhost:8000")
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatal(err)
	}
}
