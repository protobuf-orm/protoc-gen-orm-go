package apptest_test

import (
	"errors"
	"slices"
	"testing"

	"github.com/protobuf-orm/protoc-gen-orm-go/internal/apptest"
)

// The stack helpers are generated into the message package, so this is where
// they are exercised: an app that used them would be testing the same code with
// its own entity names on it.

// fake is a middleware that does nothing but say where it is in the stack.
type fake struct {
	apptest.Overlay
	name string
}

// dbHolder is a capability rather than a service: the kind of thing [Find] is
// asked for by interface, and the reason its type parameter is unconstrained.
type dbHolder interface {
	Db() string
}

type withDb struct {
	apptest.Overlay
	db string
}

func (s withDb) Db() string { return s.db }

func build(name string, err error) apptest.Builder {
	return apptest.BuilderFunc(func(next apptest.Server) (apptest.Server, error) {
		if err != nil {
			return nil, err
		}

		return fake{apptest.NewOverlay(next), name}, nil
	})
}

func names(s apptest.Server) []string {
	vs := []string{}
	for s := range apptest.Iter(s) {
		if v, ok := s.(fake); ok {
			vs = append(vs, v.name)
		}
	}

	return vs
}

func TestBuild(t *testing.T) {
	t.Run("the last builder handles the request first", func(t *testing.T) {
		sink := apptest.UnimplementedServer{}
		s, err := apptest.Build(sink, build("core", nil), build("log", nil))
		if err != nil {
			t.Fatalf("build: %s", err)
		}
		if got := names(s); !slices.Equal(got, []string{"log", "core"}) {
			t.Errorf("names = %v, want [log core]", got)
		}
		if apptest.SinkOf(s) != apptest.Server(sink) {
			t.Errorf("SinkOf is not the server Build was given")
		}
	})

	t.Run("stops at the first failure", func(t *testing.T) {
		want := errors.New("cannot be built")
		built := 0
		count := apptest.BuilderFunc(func(next apptest.Server) (apptest.Server, error) {
			built++
			return next, nil
		})

		_, err := apptest.Build(apptest.UnimplementedServer{}, count, build("log", want), count)
		if !errors.Is(err, want) {
			t.Fatalf("err = %v, want %v", err, want)
		}
		// It names the builder that failed, which is how a stack of anonymous
		// functions says which one it was.
		if err.Error() == want.Error() {
			t.Errorf("err = %q, which does not say which builder failed", err)
		}
		if built != 1 {
			t.Errorf("built = %d, want 1: it did not stop", built)
		}
	})

	t.Run("nothing to stack is the sink itself", func(t *testing.T) {
		sink := apptest.UnimplementedServer{}
		s, err := apptest.Build(sink)
		if err != nil {
			t.Fatalf("build: %s", err)
		}
		if s != apptest.Server(sink) {
			t.Errorf("got a different server")
		}
	})
}

func TestFind(t *testing.T) {
	// A stack whose innermost layer holds something, with two layers in front
	// of it that do not.
	stack := func(t *testing.T) apptest.Server {
		t.Helper()
		s, err := apptest.Build(
			withDb{apptest.NewOverlay(apptest.UnimplementedServer{}), "postgres"},
			build("core", nil),
			build("log", nil),
		)
		if err != nil {
			t.Fatalf("build: %s", err)
		}
		return s
	}

	t.Run("by the concrete type of a layer", func(t *testing.T) {
		v, ok := apptest.Find[withDb](stack(t))
		if !ok {
			t.Fatal("not found")
		}
		if v.db != "postgres" {
			t.Errorf("db = %q", v.db)
		}
	})

	// This is the one the tight constraint used to refuse: dbHolder has a
	// single method and is not a Server, so `Find[T Server]` would not compile.
	t.Run("by an interface naming what is needed", func(t *testing.T) {
		v, ok := apptest.Find[dbHolder](stack(t))
		if !ok {
			t.Fatal("not found")
		}
		if v.Db() != "postgres" {
			t.Errorf("Db = %q", v.Db())
		}
	})

	t.Run("the outermost match wins", func(t *testing.T) {
		s, err := apptest.Build(
			withDb{apptest.NewOverlay(apptest.UnimplementedServer{}), "inner"},
			func() apptest.Builder {
				return apptest.BuilderFunc(func(next apptest.Server) (apptest.Server, error) {
					return withDb{apptest.NewOverlay(next), "outer"}, nil
				})
			}(),
		)
		if err != nil {
			t.Fatalf("build: %s", err)
		}
		v, ok := apptest.Find[dbHolder](s)
		if !ok {
			t.Fatal("not found")
		}
		if v.Db() != "outer" {
			t.Errorf("Db = %q, want outer", v.Db())
		}
	})

	t.Run("nothing in the stack is one", func(t *testing.T) {
		if _, ok := apptest.Find[dbHolder](apptest.UnimplementedServer{}); ok {
			t.Error("found something")
		}
	})

	// The price of the unconstrained parameter, stated out loud: a type no
	// server could be is a legal question with a false answer rather than a
	// compile error.
	t.Run("something no server could be", func(t *testing.T) {
		if _, ok := apptest.Find[int](stack(t)); ok {
			t.Error("found an int")
		}
	})
}

func TestIter(t *testing.T) {
	t.Run("outermost first, and it ends at the sink", func(t *testing.T) {
		sink := apptest.UnimplementedServer{}
		s, err := apptest.Build(sink, build("core", nil), build("log", nil))
		if err != nil {
			t.Fatalf("build: %s", err)
		}

		n := 0
		var last apptest.Server
		for v := range apptest.Iter(s) {
			n++
			last = v
		}
		if n != 3 {
			t.Errorf("walked %d servers, want 3", n)
		}
		if last != apptest.Server(sink) {
			t.Errorf("did not end at the sink")
		}
	})

	t.Run("it stops when asked to", func(t *testing.T) {
		s, err := apptest.Build(apptest.UnimplementedServer{}, build("core", nil), build("log", nil))
		if err != nil {
			t.Fatalf("build: %s", err)
		}

		n := 0
		for range apptest.Iter(s) {
			n++
			break
		}
		if n != 1 {
			t.Errorf("walked %d servers after breaking, want 1", n)
		}
	})
}

// Both of the generated aggregates are a Server as a value, which is what lets
// a stack be asked about either of them by type.
func TestAggregatesAreServers(t *testing.T) {
	var _ apptest.Server = apptest.UnimplementedServer{}
	var _ apptest.Server = apptest.StaticServer{}
	var _ apptest.Server = &apptest.StaticServer{}

	s := apptest.StaticServer{TenantServer: apptest.UnimplementedTenantServiceServer{}}
	v, ok := apptest.Find[apptest.StaticServer](apptest.Server(s))
	if !ok {
		t.Fatal("a StaticServer sink cannot be found by its own type")
	}
	if v.TenantServer == nil {
		t.Error("found an empty one")
	}
}
