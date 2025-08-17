package app

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/protobuf-orm/protobuf-orm/graph"
	"github.com/protobuf-orm/protobuf-orm/ormpb"
	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/descriptorpb"
)

var (
	identTime = protogen.GoImportPath("time").Ident("Time")
)

func (w *fileWork) goTypeScalar(t ormpb.Type) string {
	switch t {
	case ormpb.Type_TYPE_BOOL:
		return "bool"
	case ormpb.Type_TYPE_INT32:
		return "int32"
	case ormpb.Type_TYPE_SINT32:
		return "int32"
	case ormpb.Type_TYPE_UINT32:
		return "uint32"
	case ormpb.Type_TYPE_INT64:
		return "int64"
	case ormpb.Type_TYPE_SINT64:
		return "int64"
	case ormpb.Type_TYPE_UINT64:
		return "uint64"
	case ormpb.Type_TYPE_SFIXED32:
		return "int32"
	case ormpb.Type_TYPE_FIXED32:
		return "uint32"
	case ormpb.Type_TYPE_FLOAT:
		return "float32"
	case ormpb.Type_TYPE_SFIXED64:
		return "int64"
	case ormpb.Type_TYPE_FIXED64:
		return "uint64"
	case ormpb.Type_TYPE_DOUBLE:
		return "float64"
	case ormpb.Type_TYPE_STRING:
		return "string"
	case ormpb.Type_TYPE_BYTES:
		return "[]byte"
	case ormpb.Type_TYPE_UUID:
		return "[]byte"
	case ormpb.Type_TYPE_TIME:
		return w.QualifiedGoIdent(identTime)
	}

	panic(fmt.Sprintf("must be a scalar type: %v", t.String()))
}

func (w *fileWork) goType(t ormpb.Type, s graph.Shape) string {
	if t == ormpb.Type_TYPE_GROUP {
		panic("TODO")
	}
	if t.IsScalar() {
		return w.goTypeScalar(t)
	}

	switch s := s.(type) {
	case graph.ScalarShape:
		panic("it must not be a scalar")
	case graph.MapShape:
		t := w.goType(s.V, s.S)
		return fmt.Sprintf("map[%s]%s", w.goTypeScalar(s.K), t)
	case graph.MessageShape:
		p := MustGetGoImportPath(s.Descriptor.ParentFile())
		return w.QualifiedGoIdent(p.Ident(string(s.FullName.Name())))
	default:
		panic(fmt.Sprintf("unknown shape: %s", reflect.TypeOf(s).Name()))
	}
}

func (w *fileWork) GoType(f graph.Field) string {
	return w.goType(f.Type(), f.Shape())
}

func GetGoImportPath(d protoreflect.FileDescriptor) (protogen.GoImportPath, bool) {
	opts := d.Options().(*descriptorpb.FileOptions)
	v := opts.GetGoPackage()
	if v == "" {
		return "", false
	}

	es := strings.SplitN(v, ";", 2)
	return protogen.GoImportPath(es[0]), true
}

func MustGetGoImportPath(d protoreflect.FileDescriptor) protogen.GoImportPath {
	v, ok := GetGoImportPath(d.ParentFile())
	if !ok {
		panic(fmt.Sprintf("Go import path for %s not found", d.FullName))
	}

	return v
}
