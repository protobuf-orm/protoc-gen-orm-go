package app

import (
	"fmt"
	"strings"

	"github.com/protobuf-orm/protobuf-orm/graph"
	"github.com/protobuf-orm/protoc-gen-orm-go/internal/strs"
	"google.golang.org/protobuf/compiler/protogen"
)

func (w *fileWork) xFnRefByField(p graph.Field) protogen.GoIdent {
	name_x := w.entity.Name()
	name_p := strs.GoCamelCase(p.Name())
	name := name_x + "By" + name_p

	return w.define(name, func() {
		x_ref := name_x + "Ref"

		w.P("func ", name, " (v ", w.useGoType(p.Descriptor(), p.Type().Decay()), ") ", "*"+x_ref, " {")
		w.P("	x := &", x_ref, "{}")
		w.P("	x.Set" + name_p + "(v)")
		w.P("	return x")
		w.P("}")
		w.P("")
	})
}

func (w *fileWork) xFnGetByField(p graph.Field) protogen.GoIdent {
	name_x := w.entity.Name()
	name_p := strs.GoCamelCase(p.Name())
	name := name_x + "GetBy" + name_p

	return w.define(name, func() {
		x_get := name_x + "GetRequest"

		w.P("func ", name, "(v ", w.useGoType(p.Descriptor(), p.Type().Decay()), ") ", "*"+x_get, " {")
		w.P("	return ", x_get+"_builder", "{Ref: ", w.xFnRefByField(p).GoName, "(v)}.Build()")
		w.P("}")
		w.P("")
	})
}

func (w *fileWork) xFnRefByIndex(p graph.Index) protogen.GoIdent {
	name_x := w.entity.Name()
	name_p := strs.GoCamelCase(p.Name())
	name := name_x + "By" + name_p

	return w.define(name, func() {
		x_ref := name_x + "Ref"

		args := []string{}
		for p := range p.Props() {
			t := w.useGoType(p.Descriptor(), p.Type().Decay())
			if _, ok := p.(graph.Edge); ok {
				t = "*" + t + "Ref"
			}

			name_a := p.Name()
			args = append(args, fmt.Sprintf("%s %s", name_a, t))
		}

		w.P("func ", name, "(", strings.Join(args, ", "), ") ", "*"+x_ref, " {")
		w.P("	x := &", x_ref+"By"+name_p, "{}")
		for p := range p.Props() {
			name_a := p.Name()
			w.P("	x.Set", strs.GoCamelCase(p.Name()), "(", name_a, ")")
		}
		w.P("	return ", x_ref+"_builder", "{", name_p, ": x}.Build()")
		w.P("}")
		w.P("")
	})
}

func (w *fileWork) xFnGetByIndex(p graph.Index) protogen.GoIdent {
	name_x := w.entity.Name()
	name_p := strs.GoCamelCase(p.Name())
	name := name_x + "GetBy" + name_p

	return w.define(name, func() {
		x_get := name_x + "GetRequest"

		args := []string{}
		arg_names := []string{}
		for p := range p.Props() {
			t := w.useGoTypeOf(p)
			if _, ok := p.(graph.Edge); ok {
				t = "*" + t + "Ref"
			}

			name_a := p.Name()
			args = append(args, fmt.Sprintf("%s %s", name_a, t))
			arg_names = append(arg_names, name_a)
		}

		w.P("func ", name, "(", strings.Join(args, ", "), ") ", "*"+x_get, " {")
		w.P("	return ", x_get+"_builder", "{Ref: ", w.xFnRefByIndex(p).GoName, "(", strings.Join(arg_names, ", "), ")}.Build()")
		w.P("}")
		w.P("")
	})
}
