package app

import (
	"context"
	"fmt"

	"github.com/protobuf-orm/protobuf-orm/graph"
	"google.golang.org/protobuf/compiler/protogen"
)

type work struct {
	// Qualified name of identifier
	idents map[string]bool
	// Full name of entity -> Go import path
	imports map[string]protogen.GoImportPath
}

func newWork() *work {
	return &work{
		idents:  map[string]bool{},
		imports: map[string]protogen.GoImportPath{},
	}
}

type fileWork struct {
	*protogen.GeneratedFile

	root   *work
	entity graph.Entity
}

func (w *work) newFileWork(file *protogen.GeneratedFile, entity graph.Entity) *fileWork {
	fw := &fileWork{
		GeneratedFile: file,

		root:   w,
		entity: entity,
	}

	return fw
}

func (w *fileWork) Pf(format string, a ...any) {
	fmt.Fprintf(w, format, a...)
}

func (w *work) run(ctx context.Context, gf *protogen.GeneratedFile, entity graph.Entity) error {
	fw := w.newFileWork(gf, entity)
	for p := range entity.Props() {
		if p.IsUnique() {
			fw.xFnRefByProp(p)
		}
	}

	return nil
}
