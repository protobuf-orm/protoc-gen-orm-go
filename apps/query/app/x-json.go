package app

import (
	"google.golang.org/protobuf/compiler/protogen"
)

func (w *fileWork) xRcvJson() {
	name_x := w.entity.Name()

	json := protogen.GoImportPath("google.golang.org/protobuf/encoding/protojson")
	w.P("func (x ", "*"+name_x, ") MarshalJSON() ([]byte, error) { return ", json.Ident("Marshal"), "(x) }")
	w.P("func (x ", "*"+name_x, ") UnmarshalJSON(b []byte) (error) { return ", json.Ident("Unmarshal"), "(b, x) }")
	w.P("")
}
