package app

import (
	"slices"

	"github.com/ettle/strcase"
	"github.com/protobuf-orm/protobuf-orm/graph"
	"github.com/protobuf-orm/protobuf-orm/ormpb"
	"google.golang.org/protobuf/compiler/protogen"
)

func (w *fileWork) xRcvPick() {
	n_entity := string(w.entity.FullName().Name())
	w.P("func (x ", "*"+n_entity, ") Pick() ", "*"+n_entity+"Ref", "{")
	w.P("	return ", w.xFnRefByField(w.entity.Key()), "(x.Get", strcase.ToPascal(string(w.entity.Key().FullName().Name())), "())")
	w.P("}")
	w.P("")
}

// func (x *AuditGetRequest) Picks(v *Audit) bool {
// 	switch x.WhichKey() {
// 	case AuditGetRequest_Id_case:
// 		return bytes.Equal(x.GetId(), v.GetId())

// 	default:
// 		return false
// 	}
// }

func (w *fileWork) xRcvPicks() {
	n_entity := string(w.entity.FullName().Name())
	w.P("func (x ", "*"+n_entity+"Ref) Picks(v ", "*"+n_entity, ") bool {")
	w.P("	switch x.WhichKey() {")
	for p := range w.entity.Props() {
		if !p.IsUnique() {
			continue
		}

		w.P("	case ", n_entity, "Ref_", strcase.ToPascal(string(p.FullName().Name())), "_case:")
		w.Pf("		return ")
		switch p_ := p.(type) {
		case (graph.Field):
			g := "Get" + strcase.ToPascal(string(p_.FullName().Name())) + "()"
			t := p_.Type()
			switch t {
			case
				ormpb.Type_TYPE_BYTES,
				ormpb.Type_TYPE_UUID:
				w.P(protogen.GoImportPath("bytes").Ident("Equal"), "(x.", g, ", v.", g, ")")
			default:
				w.P("x.", g, " == ", "v.", g)
			}

		case (graph.Edge):
			panic("picks for edge not implemented")
		default:
			panic("unknown type of graph prop")
		}
	}
	for p := range w.entity.Indexes() {
		if !p.IsUnique() {
			continue
		}

		np := strcase.ToPascal(p.Name())
		w.P("	case ", n_entity, "Ref_", np, "_case:")
		w.P("		x := x.Get", np, "()")
		w.Pf("		return ")

		ps := slices.Collect(p.Props())
		for i, p := range ps {
			g := "Get" + strcase.ToPascal(string(p.FullName().Name())) + "()"
			switch p.(type) {
			case (graph.Field):
				w.Pf("(x.%s == v.%s)", g, g)
			case (graph.Edge):
				w.Pf("(x.%s.Picks(v.%s))", g, g)
			default:
				panic("unknown type of graph prop")
			}
			if i+1 < len(ps) {
				w.P(" &&")
			} else {
				w.P("")
			}
		}
	}
	w.P("	default:")
	w.P("		return false")
	w.P("	}")
	w.P("}")
	w.P("")
}
