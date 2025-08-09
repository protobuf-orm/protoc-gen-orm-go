package main

import (
	"context"
	"fmt"
	"text/template"

	"github.com/protobuf-orm/protobuf-orm/graph"
	"github.com/protobuf-orm/protoc-gen-orm-go/apps/query/app"
	"google.golang.org/protobuf/compiler/protogen"
)

type QueryOpts struct {
	Namer string
}

func (h *QueryOpts) Run(ctx context.Context, p *protogen.Plugin, g *graph.Graph) error {
	opts := []app.Option{}
	if h.Namer != "" {
		v, err := template.New("namer").Parse(h.Namer)
		if err != nil {
			return fmt.Errorf("opt.query.namer: %w", err)
		}
		opts = append(opts, app.WithNamer(v))
	}

	app, err := app.New(opts...)
	if err != nil {
		return fmt.Errorf("initialize plugin: %w", err)
	}

	return app.Run(ctx, p, g)
}
