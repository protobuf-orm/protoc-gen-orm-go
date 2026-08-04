package app

import (
	"github.com/protobuf-orm/protobuf-orm/graph"
	"google.golang.org/protobuf/compiler/protogen"
)

var (
	grpc = protogen.GoImportPath("google.golang.org/grpc")

	// What the stack helpers need, and all they need. They land in the message
	// package, which every consumer imports, so it should not grow a dependency
	// on their account.
	iter   = protogen.GoImportPath("iter")
	fmtpkg = protogen.GoImportPath("fmt")
)

type Work struct {
	*protogen.GeneratedFile

	Entities []graph.Entity
	Package  protogen.GoImportPath
}
