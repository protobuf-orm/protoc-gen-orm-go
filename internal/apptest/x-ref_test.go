package apptest_test

import (
	"testing"
	"time"

	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/protobuf-orm/protoc-gen-orm-go/internal/apptest"
)

// A reference by an index with a timestamp in it. The constructor took the
// timestamp by value, so it did not compile, and `Picks` compared the two
// timestamps as pointers -- so a reference built from a row never picked that
// row unless the two happened to share one message.
func TestRefByAnIndexWithATimestamp(t *testing.T) {
	at := time.Date(2026, 10, 10, 9, 0, 0, 0, time.UTC)
	row := apptest.Shift_builder{Desk: "b", StartsAt: timestamppb.New(at)}.Build()

	// A fresh message of the same instant, which is what a reference built
	// anywhere else holds.
	if !apptest.ShiftBySlot("b", timestamppb.New(at)).Picks(row) {
		t.Error("a reference by (desk, starts_at) did not pick the row it names")
	}
	if apptest.ShiftBySlot("a", timestamppb.New(at)).Picks(row) {
		t.Error("it picked a row at another desk")
	}
	if apptest.ShiftBySlot("b", timestamppb.New(at.Add(time.Minute))).Picks(row) {
		t.Error("it picked a row at another time")
	}

	// A row with no identifier is named by the index, and its own reference
	// picks it.
	ref := row.Ref()
	if ref == nil || ref.WhichKey() != apptest.ShiftRef_Slot_case {
		t.Fatalf("a row with no identifier is not named by its index: %v", ref)
	}
	if !ref.Picks(row) {
		t.Error("a row's own reference does not pick it")
	}
	if got := apptest.ShiftGetBySlot("b", timestamppb.New(at)).GetRef(); !got.Picks(row) {
		t.Error("the Get request names a different row than the reference does")
	}
}
