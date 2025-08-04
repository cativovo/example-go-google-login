package main

import (
	"errors"
	"math"
	"net/http"
	"os"

	"github.com/gorilla/sessions"
)

const SessionName = "auth"

var errUserNotFound = errors.New("user not found")

type SessionStore struct {
	sessions.Store
}

func NewSessionStore() *SessionStore {
	store := sessions.NewFilesystemStore(os.TempDir(), []byte("goth-example"))
	// set the maxLength of the cookies stored on the disk to a larger number to prevent issues with:
	// securecookie: the value is too long
	// when using OpenID Connect , since this can contain a large amount of extra information in the id_token

	// Note, when using the FilesystemStore only the session.ID is written to a browser cookie, so this is explicit for the storage on disk
	store.MaxLength(math.MaxInt64)
	store.Options.Secure = true
	store.Options.HttpOnly = true
	store.Options.SameSite = http.SameSiteLaxMode

	return &SessionStore{
		Store: store,
	}
}

func (ss *SessionStore) GetAuthSession(r *http.Request) (*sessions.Session, error) {
	return ss.Get(r, SessionName)
}

func (ss *SessionStore) GetUser(r *http.Request) (User, error) {
	s, err := ss.GetAuthSession(r)
	if err != nil {
		return User{}, err
	}

	u, ok := s.Values["user"].(User)
	if !ok {
		return User{}, errUserNotFound
	}
	return u, nil
}
