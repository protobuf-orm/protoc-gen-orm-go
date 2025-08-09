package main

import (
	"context"
	"fmt"

	"github.com/protobuf-orm/protobuf-orm/graph"
	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/types/descriptorpb"
	"google.golang.org/protobuf/types/pluginpb"
)

type Handler struct {
	Store StoreOpts
	Query QueryOpts
}

func (h *Handler) Run(p *protogen.Plugin) error {
	p.SupportedEditionsMinimum = descriptorpb.Edition_EDITION_PROTO2
	p.SupportedEditionsMaximum = descriptorpb.Edition_EDITION_MAX
	p.SupportedFeatures = uint64(0 |
		pluginpb.CodeGeneratorResponse_FEATURE_PROTO3_OPTIONAL |
		pluginpb.CodeGeneratorResponse_FEATURE_SUPPORTS_EDITIONS,
	)

	ctx := context.Background()
	// TODO: set logger

	g := graph.NewGraph()
	for _, f := range p.Files {
		if err := graph.Parse(ctx, g, f.Desc); err != nil {
			return fmt.Errorf("parse entity at %s: %w", *f.Proto.Name, err)
		}
	}

	if err := h.Store.Run(ctx, p, g); err != nil {
		return fmt.Errorf("run store app: %w", err)
	}
	if err := h.Query.Run(ctx, p, g); err != nil {
		return fmt.Errorf("run query app: %w", err)
	}

	return nil
}
