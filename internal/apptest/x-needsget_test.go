package apptest_test

import (
	"reflect"
	"testing"

	"github.com/protobuf-orm/protoc-gen-orm-go/internal/apptest"
)

// An entity with writes and no Get has a Ref and what a Ref is enough for, and
// nothing that names the GetRequest or the Select the contract does not have.
func TestNoGetHasTheRefHelpersOnly(t *testing.T) {
	row := apptest.Counter_builder{Id: "a", Count: 3}.Build()
	if !row.Ref().Picks(row) {
		t.Error("a row's own Ref does not pick it")
	}
	if !apptest.CounterById("a").Picks(row) {
		t.Error("a Ref by its key does not pick the row")
	}
	if apptest.CounterById("b").Picks(row) {
		t.Error("a Ref by another key picks the row")
	}

	for _, v := range []any{row, apptest.CounterById("a")} {
		if _, ok := reflect.TypeOf(v).MethodByName("Pick"); ok {
			t.Errorf("%T has a Pick, which answers a GetRequest there is none of", v)
		}
	}
	if apptest.File_apptest_shift_svc_g_proto.Messages().ByName("CounterGetRequest") != nil {
		t.Error("the contract has a CounterGetRequest")
	}
}
