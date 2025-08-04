package main

import (
	"context"

	"github.com/danielgtaylor/huma/v2"
)

var humaConfig = huma.DefaultConfig("Todo", "v0.0.1")

func init() {
	humaConfig.Components.SecuritySchemes = map[string]*huma.SecurityScheme{
		"SessionAuth": {
			Type: "apiKey",
			In:   "cookie",
			Name: SessionName,
		},
	}
	humaConfig.DocsPath = ""
	humaConfig.CreateHooks = nil
}

var security = []map[string][]string{
	{"SessionAuth": {}},
}

type humaHandler[I, O any] func(ctx context.Context, input *I) (*O, error)
