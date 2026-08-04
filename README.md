# protoc-gen-orm-go

A `protoc`/`buf` plugin that generates **Go helper code** for
[protobuf-orm](https://github.com/protobuf-orm/protobuf-orm) entities: ergonomic
constructors and conversions around the generated message types, plus the gRPC
server/client wiring that ties an entity's service together.

It is meant to run *after* [protoc-gen-orm-service](../protoc-gen-orm-service)
(which emits the `.proto` services) and `protoc-gen-go` / `protoc-gen-go-grpc`
(which emit the Go message and gRPC stubs). This plugin adds the convenience
layer on top.

## What it generates

Two outputs:

### 1. Query helpers — `<name>.g.go` (one per source file)

For each entity it emits methods and constructors on the generated types:

```go
func (x *User) Ref() *UserRef                                    // value  → ref
func (x *UserRef) Pick() *UserGetRequest                         // ref    → get-request
func (x *User) Pick() *UserGetRequest                            // value  → get-request
func (x *UserRef) Picks(v *User) bool                            // does this ref match v?
func (x *UserGetRequest) WithSelect(f func(s *UserSelect)) *UserGetRequest

func (x *User) MarshalJSON() ([]byte, error)                     // via protojson
func (x *User) UnmarshalJSON(b []byte) error

// one constructor pair per key / unique index:
func UserById(v []byte) *UserRef
func UserGetById(v []byte) *UserGetRequest
func UserByAlias(alias string, tenant *TenantRef) *UserRef
func UserGetByAlias(alias string, tenant *TenantRef) *UserGetRequest
```

### 2. Store wiring — `store.g.go` (one per module)

Aggregates every entity's gRPC service into a single registration surface:

```go
type Server interface { User() UserServiceServer; Tenant() TenantServiceServer; ... }
func RegisterServer(g *grpc.Server, s Server)
type UnimplementedServer struct{ ... }
type StaticServer struct{ UserServer ...; TenantServer ...; ... }

type Client interface { User() UserServiceClient; ... }
func NewClient(c *grpc.ClientConn) Client
```

### 3. The stack — also `store.g.go`

A `Server` is meant to be wrapped: one implementation runs the queries and every
other one sits in front of it to add a behaviour of its own. The same handful of
pieces makes that work, and they are emitted here because every one of them is
spelled in terms of the generated `Server` — an app that wrote them itself would
write this file.

```go
type Middleware interface { Next() Server }
func Iter(s Server) iter.Seq[Server]
func Find[T any](s Server) (T, bool)
func SinkOf(s Server) Server

type Overlay struct{ Server }          // embed to implement only what you override
func NewOverlay(next Server) Overlay

type Builder interface { Build(next Server) (Server, error) }
type BuilderFunc func(next Server) (Server, error)
func Build(sink Server, mws ...Builder) (Server, error)
```

`Overlay` is why these cannot live in a shared module instead: it has to embed
`Server` to promote the services a layer does not override, and Go does not
allow embedding a type parameter. Naming `Server` means living where `Server`
lives. They add `iter` and `fmt` to the package's imports and nothing else.

**`Find` takes any type on purpose.** A stack is asked two different questions —
"which layer is the bare server", naming a concrete type, and "who here has a
database", naming a one-method interface — and a constraint of `Server` would
admit only the first, since a one-method interface does not implement `Server`.
That is what lets a capability be *found* rather than declared: a layer holding
something the others need does not have to put a method on `Server`, which would
force every layer, every `Overlay` and every helper above to be rewritten to
match. The price is that a `T` no server could ever be is a legal question with
a false answer.

## Usage

```yaml
version: v2
plugins:
  - local: [go, run, github.com/protobuf-orm/protoc-gen-orm-go]
    out: .
    opt:
      - module=github.com/your/module
      - query.namer={{ .Name }}.g.go
```

Options:

| Option         | Default              | Meaning                                          |
| -------------- | -------------------- | ------------------------------------------------ |
| `store.name`   | `store.g.go`         | Output filename of the server/client wiring.     |
| `query.namer`  | `{{ .Name }}.g.go`   | Go text/template for the per-file query helpers. |

## Structure

```
main.go                 flag parsing; wires protogen → Handler
handler.go              parses files into a graph.Graph, runs the Store + Query apps
apps/store/app/         generates store.g.go (Server/Client/registration)
apps/query/app/
  app.go                per-file driver
  work.go               type/import bookkeeping (useGoType, useGoTypeOf)
  x-ref.go              Ref()/By<Key>() constructors  (uses Type().Decay())
  x-pick.go             Pick()/Picks()/Select helpers (uses Type().Decay())
  x-select.go           WithSelect fluent builder
  x-json.go             Marshal/UnmarshalJSON
internal/strs/          protobuf name-casing helpers (vendored from protobuf-go)
```

Type mapping goes through the library's `graph.GoType`/`GoTypeOf` and
`Type.Decay()`. Note `Decay()` now folds `TYPE_ENUM` into the integer category,
so an enum-typed ref field generates a `v != 0` presence check rather than the
(invalid) `v != nil`.

## Development

```sh
buf generate     # regenerate apptest fixtures
go build ./...
go test ./...
```
