package main

import (
	"errors"
	"net/http"

	"github.com/gorilla/sessions"
)

const authSession = "auth"

var errUserNotFound = errors.New("user not found")

type SessionStore struct {
	sessions.Store
}

func NewSession(s sessions.Store) *SessionStore {
	return &SessionStore{
		Store: s,
	}
}

func (ss *SessionStore) GetAuthSession(r *http.Request) (*sessions.Session, error) {
	return ss.Get(r, authSession)
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
