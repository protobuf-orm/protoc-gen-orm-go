package app

import (
	"github.com/ettle/strcase"
	"github.com/protobuf-orm/protobuf-orm/graph"
	"google.golang.org/protobuf/compiler/protogen"
)

func (w *fileWork) xFnRefByProp(p graph.Prop) protogen.GoIdent {
	switch p_ := p.(type) {
	case (graph.Field):
		return w.xFnRefByField(p_)
	case (graph.Index):
		panic("not implemented")
		// return w.xFnRefByIndex(p_)
	default:
		panic("unknown type of graph prop")
	}
}

func (w *fileWork) xFnRefByField(f graph.Field) protogen.GoIdent {
	n_entity := string(w.entity.FullName().Name())
	n_field := string(f.FullName().Name())
	name := n_entity + "By" + strcase.ToPascal(n_field)

	n_ref := n_entity + "Ref"

	gt := w.GoType(f)
	w.P("func ", name, " (v ", gt, ") ", "*"+n_ref, " {")
	w.P("	x := ", "&"+n_ref+"{}")
	w.P("	x.Set" + strcase.ToPascal(n_field) + "(v)")
	w.P("	return x")
	w.P("}")

	p, ok := w.root.imports[string(w.entity.FullName())]
	if !ok {
		panic("import path for the entity must be exist")
	}

	return p.Ident(name)
}

// func (w *fileWork) xFnRefByIndex(i graph.Index) protogen.GoIdent {

// }
