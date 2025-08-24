package app

import (
	"context"
	"fmt"

	"github.com/protobuf-orm/protobuf-orm/graph"
	"github.com/protobuf-orm/protobuf-orm/ormpb"
	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/reflect/protoreflect"
)

type work struct {
	// Qualified name of the identifier defined by this app.
	idents map[protogen.GoIdent]bool
}

func newWork() *work {
	return &work{
		idents: map[protogen.GoIdent]bool{},
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
	return &fileWork{
		GeneratedFile: file,

		root:   w,
		entity: entity,
		pkg:    graph.MustGetGoImportPath(entity.Descriptor().ParentFile()),

		deferred: []func(){},
	}
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

func (w *fileWork) useGoTypeOf(p graph.Prop) string {
	return graph.GoTypeOf(p, w.QualifiedGoIdent)
}

func (w *fileWork) useGoType(d protoreflect.FieldDescriptor, t ormpb.Type) string {
	return graph.GoType(d, t, w.QualifiedGoIdent)
}

func (w *work) run(ctx context.Context, gf *protogen.GeneratedFile, entity graph.Entity) error {
	fw := w.newFileWork(gf, entity)
	fw.xRcvPick()
	fw.xRcvPickUp()
	fw.xRcvPicks()
	fw.xRcvWithSelect()
	fw.xRcvJson()
	for p := range entity.Keys() {
		switch p := p.(type) {
		case graph.Field:
			fw.xFnRefByField(p)
			fw.xFnGetByField(p)

		case graph.Edge:
			panic("ref for edge not implemented")

		case graph.Index:
			fw.xFnRefByIndex(p)
			fw.xFnGetByIndex(p)

		default:
			panic("unknown type of graph prop")
		}
	}

	for _, f := range fw.deferred {
		f()
	}
	return nil
}
