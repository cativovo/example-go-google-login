package main

import (
	"fmt"
	"html/template"
	"net/http"

	"github.com/markbates/goth/gothic"
)

func authCallback(store *SessionStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
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

		s, err := store.GetAuthSession(r)
		if err != nil {
			fmt.Fprintln(w, fmt.Errorf("dito get: %w", err))
			return
		}

		u := User{
			ID:        gothicUser.UserID,
			FirstName: gothicUser.FirstName,
			LastName:  gothicUser.LastName,
			Email:     gothicUser.Email,
			AvatarURL: gothicUser.AvatarURL,
		}

		s.Values["user"] = u
		s.Options.MaxAge = int(sessionDuration.Seconds())

		if err := s.Save(r, w); err != nil {
			fmt.Fprintln(w, fmt.Errorf("dito save: %w", err))
			return
		}

		w.Header().Set("Location", "/")
		w.WriteHeader(http.StatusTemporaryRedirect)
	}
}

func logout(store *SessionStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		gothic.Logout(w, r)
		s, err := store.GetAuthSession(r)
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
	}
}

func login() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		t, _ := template.New("foo").Parse(loginTemplate)
		t.Execute(w, nil)
	}
}

func getUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		u := userFromContext(r.Context())
		t, _ := template.New("foo").Parse(userTemplate)
		t.Execute(w, map[string]any{
			"FirstName": u.FirstName,
			"LastName":  u.LastName,
			"Email":     u.Email,
			"ID":        u.ID,
			"AvatarURL": u.AvatarURL,
		})
	}
}

var loginTemplate = `<p><a href="/auth/begin">Log in with google</a></p>`

var userTemplate = `
<p><a href="/auth/logout">logout</a></p>
<p>Name: [{{.LastName}}, {{.FirstName}}]</p>
<p>Email: {{.Email}}</p>
<p>AvatarURL: {{.AvatarURL}} <img src="{{.AvatarURL}}"></p>
<p>UserID: {{.ID}}</p>
`
