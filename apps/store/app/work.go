package app

import (
	"github.com/protobuf-orm/protobuf-orm/graph"
	"google.golang.org/protobuf/compiler/protogen"
)

var (
	grpc = protogen.GoImportPath("google.golang.org/grpc")
)

type Work struct {
	*protogen.GeneratedFile

	Entities []graph.Entity
	Package  protogen.GoImportPath
}
