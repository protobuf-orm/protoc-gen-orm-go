package app

import (
	"fmt"
	"slices"
	"strings"

	"github.com/protobuf-orm/protobuf-orm/graph"
	"github.com/protobuf-orm/protobuf-orm/ormpb"
	"github.com/protobuf-orm/protoc-gen-orm-go/internal/strs"
	"google.golang.org/protobuf/compiler/protogen"
)

func (w *fileWork) xRefRcvPick() {
	name_x := w.entity.Name()
	x_get := name_x + "GetRequest"

	w.P("func (x ", "*"+name_x+"Ref", ") Pick() ", "*"+x_get, "{")
	w.P("	return ", x_get+"_builder", "{Ref: x}.Build()")
	w.P("}")
	w.P("")
}

func (w *fileWork) xRcvRef() {
	name_x := w.entity.Name()

	w.P("func (x ", "*"+name_x, ") Ref() ", "*"+name_x+"Ref", " {")
	for p := range w.entity.Keys() {
		name_p := strs.GoCamelCase(p.Name())
		u := "x.Get" + name_p + "()"

		switch p := p.(type) {
		case (graph.Field):
			c := ""
			switch p.Type().Decay() {
			case ormpb.Type_TYPE_MESSAGE:
				c = "v != nil"
			case ormpb.Type_TYPE_BYTES, ormpb.Type_TYPE_STRING:
				c = "len(v) > 0"
			case ormpb.Type_TYPE_FLOAT, ormpb.Type_TYPE_INT, ormpb.Type_TYPE_UINT:
				c = "v != 0"
			}
			w.P("	if v := ", u, "; ", c, " {")
			w.P("		return ", w.xFnRefByField(p), "(v)")
			w.P("	}")

		case (graph.Edge):
			panic("pick for edge not implemented")

		case (graph.Index):
			w.P("	{")
			i := 0
			for p := range p.Props() {
				i++
				v := fmt.Sprintf("v%d", i)
				w.P("		", v, " := ", "x.Get"+strs.GoCamelCase(p.Name()), "()")
			}

			i = 0
			cs := []string{}
			vs := []string{}
			for p := range p.Props() {
				i++
				c := ""
				v := fmt.Sprintf("v%d", i)

				switch p.Type().Decay() {
				case ormpb.Type_TYPE_MESSAGE:
					c = fmt.Sprintf("%s != nil", v)
				case ormpb.Type_TYPE_BYTES, ormpb.Type_TYPE_STRING:
					c = fmt.Sprintf("len(%s) > 0", v)
				case ormpb.Type_TYPE_FLOAT, ormpb.Type_TYPE_INT, ormpb.Type_TYPE_UINT:
					c = fmt.Sprintf("%s != 0", v)
				}

				switch p.(type) {
				case (graph.Field):
				case (graph.Edge):
					v = v + ".Ref()"
				default:
					panic("unknown type of index prop")
				}
				cs = append(cs, c)
				vs = append(vs, v)
			}
			w.P("		if ", strings.Join(cs, " && "), " {")
			w.P("			return ", name_x+"By"+name_p, "(", strings.Join(vs, ", "), ")")
			w.P("		}")
			w.P("	}")

		default:
			panic("unknown type of graph element")
		}
	}
	w.P("")
	w.P("	return nil")
	w.P("}")
	w.P("")
}

func (w *fileWork) xRcvPick() {
	name_x := w.entity.Name()
	x_get := name_x + "GetRequest"

	w.P("func (x ", "*"+name_x, ") Pick() ", "*"+x_get, "{")
	w.P("	return x.Ref().Pick()")
	w.P("}")
	w.P("")
}

func (w *fileWork) xRcvPicks() {
	name_x := w.entity.Name()
	x_ref := name_x + "Ref"

	w.P("func (x ", "*"+x_ref, ") Picks(v ", "*"+name_x, ") bool {")
	w.P("	switch x.WhichKey() {")
	for p := range w.entity.Keys() {
		name_p := strs.GoCamelCase(p.Name())
		u := "Get" + name_p + "()"

		w.P("	case ", x_ref+"_"+name_p+"_case", ":")
		switch p := p.(type) {
		case (graph.Field):
			a := "x." + u
			b := "v." + u
			switch p.Type().Decay() {
			case ormpb.Type_TYPE_BYTES:
				w.P("		return ", protogen.GoImportPath("bytes").Ident("Equal"), "(", a, ", ", b, ")")
			default:
				w.P("		return ", a, " == ", b)
			}

		case (graph.Edge):
			panic("pick for edge not implemented")

		case (graph.Index):
			w.P("		x := x." + u)
			w.Pf("		return ")

			ps := slices.Collect(p.Props())
			for i, p := range ps {
				name_p := strs.GoCamelCase(p.Name())
				u := "Get" + name_p + "()"
				switch p.(type) {
				case (graph.Field):
					w.Pf("(x.%s == v.%s)", u, u)
				case (graph.Edge):
					w.Pf("(x.%s.Picks(v.%s))", u, u)
				default:
					panic("unknown type of graph prop")
				}
				if i+1 < len(ps) {
					w.P(" &&")
				} else {
					w.P("")
				}
			}

		default:
			panic("unknown type of graph element")
		}
	}
	w.P("	default:")
	w.P("		return false")
	w.P("	}")
	w.P("}")
	w.P("")
}
