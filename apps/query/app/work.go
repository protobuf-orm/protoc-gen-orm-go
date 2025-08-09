package app

import (
	"context"
	"fmt"

	"github.com/protobuf-orm/protobuf-orm/graph"
	"google.golang.org/protobuf/compiler/protogen"
)

type work struct {
	// Qualified name of identifier
	idents map[protogen.GoIdent]bool
	// Full name of entity -> Go import path
	imports map[string]protogen.GoImportPath
}

func newWork() *work {
	return &work{
		idents:  map[protogen.GoIdent]bool{},
		imports: map[string]protogen.GoImportPath{},
	}
}

type fileWork struct {
	*protogen.GeneratedFile

	root   *work
	entity graph.Entity
	pkg    protogen.GoImportPath

	deferred []func()
}

func (w *work) newFileWork(file *protogen.GeneratedFile, entity graph.Entity) *fileWork {
	pkg, ok := w.imports[string(entity.FullName())]
	if !ok {
		panic("import path for the entity must be exist")
	}

	fw := &fileWork{
		GeneratedFile: file,

		root:   w,
		entity: entity,
		pkg:    pkg,

		deferred: []func(){},
	}

	return fw
}

func (w *fileWork) Pf(format string, a ...any) {
	fmt.Fprintf(w, format, a...)
}

func (w *fileWork) define(name string, f func()) protogen.GoIdent {
	ident := w.pkg.Ident(name)
	if _, ok := w.root.idents[ident]; ok {
		return ident
	}

	w.deferred = append(w.deferred, f)
	w.root.idents[ident] = true
	return ident
}

func (w *work) run(ctx context.Context, gf *protogen.GeneratedFile, entity graph.Entity) error {
	fw := w.newFileWork(gf, entity)
	fw.xRcvPick()
	fw.xRcvPicks()
	for p := range entity.Props() {
		if !p.IsUnique() {
			continue
		}

		fw.xFnRefByProp(p)
	}
	for p := range entity.Indexes() {
		if !p.IsUnique() {
			continue
		}

		fw.xFnRefByIndex(p)
	}

	for _, f := range fw.deferred {
		f()
	}

	return nil
}
