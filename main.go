package main

import (
	"context"
	"encoding/gob"
	"log"
	"math"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/gorilla/sessions"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/google"
)

type User struct {
	ID        string
	FirstName string
	LastName  string
	Email     string
	AvatarURL string
}

var store *sessions.FilesystemStore

var sessionDuration = 7 * 24 * time.Hour

func init() {
	// fix for -> securecookie: error - caused by: securecookie: error - caused by: gob: type not registered for interface: main.User
	gob.Register(User{})

	store = sessions.NewFilesystemStore("/tmp", []byte("goth-example"))

	// set the maxLength of the cookies stored on the disk to a larger number to prevent issues with:
	// securecookie: the value is too long
	// when using OpenID Connect , since this can contain a large amount of extra information in the id_token

	// Note, when using the FilesystemStore only the session.ID is written to a browser cookie, so this is explicit for the storage on disk
	store.MaxLength(math.MaxInt64)
	store.Options.Secure = true
	store.Options.HttpOnly = true
	store.Options.SameSite = http.SameSiteLaxMode

	gothic.Store = store
	gothic.GetProviderName = func(_ *http.Request) (string, error) {
		return "google", nil
	}
}

func main() {
	scopes := []string{"profile", "email"}
	goth.UseProviders(
		google.New(os.Getenv("GOOGLE_KEY"), os.Getenv("GOOGLE_SECRET"), "http://127.0.0.1:8000/auth/callback", scopes...),
	)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sessionStore := NewSession(store)
	todoStore := NewTodoStore()
	r := chi.NewRouter()
	mountRoutes(r, sessionStore, todoStore)

	server := http.Server{
		Addr:    ":8000",
		Handler: r,
	}

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
