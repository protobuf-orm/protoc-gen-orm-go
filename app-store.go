package main

import (
	"context"
	"fmt"

	"github.com/protobuf-orm/protobuf-orm/graph"
	"github.com/protobuf-orm/protoc-gen-orm-go/apps/store/app"
	"google.golang.org/protobuf/compiler/protogen"
)

type StoreOpts struct {
	Name string
}

func (h *StoreOpts) Run(ctx context.Context, p *protogen.Plugin, g *graph.Graph) error {
	opts := []app.Option{}
	if h.Name != "" {
		opts = append(opts, app.WithName(h.Name))
	}

	app, err := app.New(opts...)
	if err != nil {
		return fmt.Errorf("initialize plugin: %w", err)
	}

	return app.Run(ctx, p, g)
}
