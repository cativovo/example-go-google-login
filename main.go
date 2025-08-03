package main

import (
	"encoding/gob"
	"log"
	"math"
	"net/http"
	"os"
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

	store = sessions.NewFilesystemStore(os.TempDir(), []byte("goth-example"))

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

	sessionStore := NewSession(store)
	r := chi.NewRouter()
	mountRoutes(r, sessionStore)

	log.Println("listening on localhost:8000")
	log.Fatal(http.ListenAndServe(":8000", r))
}
