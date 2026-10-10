package app

import (
	"fmt"
	"strings"

	"github.com/protobuf-orm/protobuf-orm/graph"
	"github.com/protobuf-orm/protoc-gen-orm-go/internal/strs"
	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/reflect/protoreflect"
)

// argType is the type a constructor takes for one prop of a reference: what
// the reference's setter takes, since that is where the argument goes.
//
// `useGoType` with the decayed type answers with the wire's type -- `[]byte`
// for a uuid -- and names a message without its pointer, so a message-typed
// field is given one here, the way an edge is given `*…Ref`. A timestamp in an
// index was the first to reach this, and the constructor did not compile.
func (w *fileWork) argType(p graph.Prop) string {
	t := w.useGoType(p.Descriptor(), p.Type().Decay())
	if _, ok := p.(graph.Edge); ok {
		return "*" + t + "Ref"
	}
	if d := p.Descriptor(); d.Kind() == protoreflect.MessageKind && !d.IsMap() {
		return "*" + t
	}

	return t
}

func (w *fileWork) xFnRefByField(p graph.Field) protogen.GoIdent {
	name_x := w.entity.Name()
	name_p := strs.GoCamelCase(p.Name())
	name := name_x + "By" + name_p

	return w.define(name, func() {
		x_ref := name_x + "Ref"

		w.P("func ", name, " (v ", w.argType(p), ") ", "*"+x_ref, " {")
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

		w.P("func ", name, "(v ", w.argType(p), ") ", "*"+x_get, " {")
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
			name_a := p.Name()
			args = append(args, fmt.Sprintf("%s %s", name_a, w.argType(p)))
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
			// The same type the constructor this calls takes; see
			// [fileWork.argType]. `useGoTypeOf` answers with what the column
			// is -- `uuid.UUID` for a uuid -- and the constructor takes what
			// the setter takes, which is `[]byte`. The two disagreeing meant
			// this function did not compile, which only an index holding a
			// uuid could reach.
			name_a := p.Name()
			args = append(args, fmt.Sprintf("%s %s", name_a, w.argType(p)))
			arg_names = append(arg_names, name_a)
		}

		w.P("func ", name, "(", strings.Join(args, ", "), ") ", "*"+x_get, " {")
		w.P("	return ", x_get+"_builder", "{Ref: ", w.xFnRefByIndex(p).GoName, "(", strings.Join(arg_names, ", "), ")}.Build()")
		w.P("}")
		w.P("")
	})
}
