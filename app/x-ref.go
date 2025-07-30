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
	case (graph.Edge):
		panic("ref for edge not implemented")
	default:
		panic("unknown type of graph prop")
	}
}

func (w *fileWork) xFnRefByField(p graph.Field) protogen.GoIdent {
	n_entity := string(w.entity.FullName().Name())
	n_field := string(p.FullName().Name())
	name := n_entity + "By" + strcase.ToPascal(n_field)

	return w.define(name, func() {
		t := w.GoType(p)
		n_ref := n_entity + "Ref"

		w.P("func ", name, " (v ", t, ") ", "*"+n_ref, " {")
		w.P("	x := ", "&"+n_ref+"{}")
		w.P("	x.Set" + strcase.ToPascal(n_field) + "(v)")
		w.P("	return x")
		w.P("}")
		w.P("")
	})
}

func (w *fileWork) xFnRefByIndex(p graph.Index) protogen.GoIdent {
	n_entity := string(w.entity.FullName().Name())
	n_field := strcase.ToPascal(string(p.Name()))
	name := n_entity + "By" + n_field

	n_ref := n_entity + "RefBy" + n_field

	w.Pf("func %s (", name)
	for p := range p.Props() {
		t := ""
		switch p_ := p.(type) {
		case graph.Field:
			t = w.GoType(p_)
		case graph.Edge:
			target := p_.Target()
			pkg, ok := w.root.imports[string(target.FullName())]
			if !ok {
				panic("import path for the entity must be exist")
			}

			n_target := string(target.FullName().Name())
			t = "*" + w.QualifiedGoIdent(pkg.Ident(n_target+"Ref"))
		default:
			panic("unknown type of graph prop")
		}

		w.Pf("%s %s,", p.FullName().Name(), t)
	}
	w.P(") *"+n_ref, " {")
	w.P("	x := ", "&"+n_ref+"{}")
	for p := range p.Props() {
		n_field := string(p.FullName().Name())
		w.P("	x.Set", strcase.ToPascal(n_field), "(", n_field, ")")
	}
	w.P("	return x")
	w.P("}")
	w.P("")

	pkg, ok := w.root.imports[string(w.entity.FullName())]
	if !ok {
		panic("import path for the entity must be exist")
	}

	return pkg.Ident(name)
}
