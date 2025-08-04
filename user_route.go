package main

import (
	"net/http"

	"github.com/danielgtaylor/huma/v2"
)

func mountUserRoutes(api huma.API, sessionStore *SessionStore) {
	userGroup := huma.NewGroup(api, "/user")
	huma.Register(
		userGroup,
		huma.Operation{
			OperationID: "get-user",
			Method:      http.MethodGet,
			Path:        "",
			Summary:     "Get the currently signed in user",
			Security:    security,
		},
		getUser(sessionStore),
	)
}
