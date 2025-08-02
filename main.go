package main

import (
	"context"
	"encoding/gob"
	"fmt"
	"html/template"
	"log"
	"math"
	"net/http"
	"os"
	"sort"
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

var sessionExpiration = 7 * 24 * time.Hour

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
}

func auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		s, err := store.Get(r, "test")
		if err != nil {
			fmt.Fprintf(w, "here: %s", err)
			return
		}

		user, ok := s.Values["user"].(User)
		if !ok {
			http.Redirect(w, r, "/login", http.StatusFound)
			return
		}

		r = r.WithContext(context.WithValue(r.Context(), "user", user))
		next.ServeHTTP(w, r)
	})
}

func withProvider(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r = gothic.GetContextWithProvider(r, "google")
		next.ServeHTTP(w, r)
	})
}

func main() {
	scopes := []string{"profile", "email"}
	goth.UseProviders(
		google.New(os.Getenv("GOOGLE_KEY"), os.Getenv("GOOGLE_SECRET"), "http://127.0.0.1:8000/auth/google/callback", scopes...),
	)

	m := map[string]string{
		"google": "Google",
	}
	var keys []string
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	providerIndex := &ProviderIndex{Providers: keys, ProvidersMap: m}

	r := chi.NewRouter()
	r.Use(withProvider)
	r.Get("/auth/{provider}/callback", func(w http.ResponseWriter, r *http.Request) {
		gothicUser, err := gothic.CompleteUserAuth(w, r)
		if err != nil {
			fmt.Fprintln(w, fmt.Errorf("dito complete: %w", err))
			return
		}

		if v := gothicUser.RawData["verified_email"].(bool); !v {
			w.WriteHeader(http.StatusUnauthorized)
			fmt.Fprintln(w, "email not verified")
			return
		}

		s, err := store.Get(r, "test")
		if err != nil {
			fmt.Fprintln(w, fmt.Errorf("dito get: %w", err))
			return
		}

		user := User{
			ID:        gothicUser.UserID,
			FirstName: gothicUser.FirstName,
			LastName:  gothicUser.LastName,
			Email:     gothicUser.Email,
			AvatarURL: gothicUser.AvatarURL,
		}

		s.Values["user"] = user
		s.Options.MaxAge = int(sessionExpiration.Seconds())

		if err := s.Save(r, w); err != nil {
			fmt.Fprintln(w, fmt.Errorf("dito save: %w", err))
			return
		}

		w.Header().Set("Location", "/")
		w.WriteHeader(http.StatusTemporaryRedirect)
	})

	r.Get("/logout/{provider}", func(w http.ResponseWriter, r *http.Request) {
		gothic.Logout(w, r)
		s, err := store.Get(r, "test")
		if err != nil {
			fmt.Fprintln(w, err)
			return
		}

		s.Options.MaxAge = -1

		if err := store.Save(r, w, s); err != nil {
			fmt.Fprintln(w, err)
			return
		}
		w.Header().Set("Location", "/")
		w.WriteHeader(http.StatusTemporaryRedirect)
	})

	r.Get("/auth/google", gothic.BeginAuthHandler)

	r.Group(func(r chi.Router) {
		r.Use(auth)
		r.Get("/", func(w http.ResponseWriter, r *http.Request) {
			user := r.Context().Value("user").(User)
			t, _ := template.New("foo").Parse(userTemplate)
			t.Execute(w, map[string]any{
				"Provider":  "google",
				"FirstName": user.FirstName,
				"LastName":  user.LastName,
				"Email":     user.Email,
				"ID":        user.ID,
				"AvatarURL": user.AvatarURL,
			})
		})
	})

	r.Get("/login", func(w http.ResponseWriter, r *http.Request) {
		t, _ := template.New("foo").Parse(indexTemplate)
		t.Execute(w, providerIndex)
	})

	log.Println("listening on localhost:8000")
	log.Fatal(http.ListenAndServe(":8000", r))
}

type ProviderIndex struct {
	Providers    []string
	ProvidersMap map[string]string
}

var indexTemplate = `{{range $key,$value:=.Providers}}
    <p><a href="/auth/{{$value}}">Log in with {{index $.ProvidersMap $value}}</a></p>
{{end}}`

var userTemplate = `
<p><a href="/logout/{{.Provider}}">logout</a></p>
<p>Name: [{{.LastName}}, {{.FirstName}}]</p>
<p>Email: {{.Email}}</p>
<p>AvatarURL: {{.AvatarURL}} <img src="{{.AvatarURL}}"></p>
<p>UserID: {{.ID}}</p>
`
