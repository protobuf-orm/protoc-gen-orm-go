package app

import "google.golang.org/protobuf/compiler/protogen"

// xPackageDoc explains the stacking model on the package itself.
//
// It is written here because the package it describes is generated: an app that
// documented it would be documenting somebody else's types, and would say it
// again in the next app.
func xPackageDoc(gf *protogen.GeneratedFile, name string) {
	gf.P("// Package ", name, " holds the messages of the schema, the service servers")
	gf.P("// generated from them, and the pieces every implementation of those servers")
	gf.P("// shares.")
	gf.P("//")
	gf.P("// A server is a [Server], which is nothing but a set of those service")
	gf.P("// servers. Implementations are stacked on top of each other: one of them")
	gf.P("// runs the queries against the database and every other one wraps it to add")
	gf.P("// a behaviour of its own, so a request walks the stack from the outermost")
	gf.P("// server down to the one that answers it. [Build] does the stacking,")
	gf.P("// [Overlay] is what a layer embeds to implement only the services it cares")
	gf.P("// about, and [Iter], [Find] and [SinkOf] look into a stack that was built.")
	gf.P("//")
	gf.P("// Whatever a layer can do besides answering a service is found with [Find]")
	gf.P("// rather than added to [Server]. That is what keeps [Server] the generated")
	gf.P("// set it is: a layer that holds a database, or a connection, or anything")
	gf.P("// else the others need, is asked for by what it can do, and every other")
	gf.P("// layer stays as narrow as it was written.")
}

// xStack emits the pieces every server implementation shares.
//
// They are generated rather than written by hand because every one of them is
// expressed in terms of [Server], which is generated: an app that wrote them
// itself would write the same file, and did -- twice, in the two apps this was
// taken from, drifting apart in the details.
//
// They are in this package and not one of the app's own for the same reason.
// [Overlay] has to embed Server to promote the services it does not override,
// and a type parameter cannot be embedded in Go, so a shared generic version is
// not expressible. Naming Server means living where Server lives.
func (w *Work) xStack() {
	w.xStackWalk()
	w.xStackOverlay()
	w.xStackBuilder()
}

// xStackWalk emits the three ways to look into a stack.
func (w *Work) xStackWalk() {
	w.P("// Middleware is a server that delegates to another server.")
	w.P("type Middleware interface {")
	w.P("	Next() Server")
	w.P("}")
	w.P("")

	w.P("// Iter yields `s` and, as long as they are middlewares, the servers behind")
	w.P("// it, from the outermost one to the one that handles the request.")
	w.P("func Iter(s Server) ", iter.Ident("Seq"), "[Server] {")
	w.P("	return func(yield func(Server) bool) {")
	w.P("		for s != nil {")
	w.P("			if !yield(s) {")
	w.P("				return")
	w.P("			}")
	w.P("")
	w.P("			mw, ok := s.(Middleware)")
	w.P("			if !ok {")
	w.P("				return")
	w.P("			}")
	w.P("")
	w.P("			s = mw.Next()")
	w.P("		}")
	w.P("	}")
	w.P("}")
	w.P("")

	w.P("// Find returns the outermost server in the stack that is a `T`.")
	w.P("//")
	w.P("// It is how a stack answers for something that is not a service. A layer")
	w.P("// that owns a database, or a connection, or anything else the others need")
	w.P("// is reached by asking the stack for what can do it, rather than by every")
	w.P("// layer carrying it or by [Server] growing a method for it.")
	w.P("//")
	w.P("// `T` is unconstrained on purpose. The question is asked two ways -- by the")
	w.P("// concrete type of a layer, and by an interface naming the one thing the")
	w.P("// caller needs -- and a constraint of Server would admit only the first,")
	w.P("// since an interface with one method on it does not implement Server. The")
	w.P("// cost is that a `T` no server could ever be is a legal question with a")
	w.P("// false answer.")
	w.P("func Find[T any](s Server) (T, bool) {")
	w.P("	for s := range Iter(s) {")
	w.P("		if v, ok := s.(T); ok {")
	w.P("			return v, true")
	w.P("		}")
	w.P("	}")
	w.P("")
	w.P("	var zero T")
	w.P("	return zero, false")
	w.P("}")
	w.P("")

	w.P("// SinkOf returns the server the stack ends at, which is the one the others")
	w.P("// were built in front of and the only one that answers out of a database")
	w.P("// rather than by asking somebody else.")
	w.P("//")
	w.P("// It is the same server [Build] was given, and it is named the same way.")
	w.P("// What is at the end is not always what a caller means, though: reach for")
	w.P("// [Find] when there is a particular server in mind, and keep this for when")
	w.P("// the stack itself is the subject.")
	w.P("func SinkOf(s Server) Server {")
	w.P("	for v := range Iter(s) {")
	w.P("		s = v")
	w.P("	}")
	w.P("	return s")
	w.P("}")
	w.P("")
}

// xStackOverlay emits the middleware base every layer embeds.
func (w *Work) xStackOverlay() {
	w.P("// Overlay is a middleware that forwards every service it does not override")
	w.P("// to the next server. Embed it to implement only the services of interest:")
	w.P("//")
	w.P("//	type Server struct {")
	w.P("//		", w.Package.Ident("Overlay"))
	w.P("//	}")
	w.P("//")
	w.P("//	func (s Server) ", w.Entities[0].Name(), "() ", w.Package.Ident(w.Entities[0].Name()+"ServiceServer"), " { ... }")
	w.P("type Overlay struct {")
	w.P("	Server")
	w.P("}")
	w.P("")
	w.P("func NewOverlay(next Server) Overlay {")
	w.P("	return Overlay{next}")
	w.P("}")
	w.P("")
	w.P("func (s Overlay) Next() Server {")
	w.P("	return s.Server")
	w.P("}")
	w.P("")
}

// xStackBuilder emits the stacking itself.
func (w *Work) xStackBuilder() {
	w.P("// Builder makes a server that sits in front of `next`. Building may fail")
	w.P("// since a server is free to open the resources it needs, or to reject the")
	w.P("// settings it was given, while it is made.")
	w.P("type Builder interface {")
	w.P("	Build(next Server) (Server, error)")
	w.P("}")
	w.P("")
	w.P("type BuilderFunc func(next Server) (Server, error)")
	w.P("")
	w.P("func (f BuilderFunc) Build(next Server) (Server, error) {")
	w.P("	return f(next)")
	w.P("}")
	w.P("")

	w.P("// Build stacks the given middlewares on top of `sink`. Each builder wraps")
	w.P("// the result of the previous one, so the last builder given handles the")
	w.P("// request first:")
	w.P("//")
	w.P("//	Build(sink, core.Build(), log.Build())")
	w.P("//	// log -> core -> sink")
	w.P("//")
	w.P("// It stops at the first builder that fails, reporting it by its type.")
	w.P("func Build(sink Server, mws ...Builder) (Server, error) {")
	w.P("	s := sink")
	w.P("	for _, mw := range mws {")
	w.P("		v, err := mw.Build(s)")
	w.P("		if err != nil {")
	w.P("			return nil, ", fmtpkg.Ident("Errorf"), "(\"%T: %w\", mw, err)")
	w.P("		}")
	w.P("")
	w.P("		s = v")
	w.P("	}")
	w.P("")
	w.P("	return s, nil")
	w.P("}")
	w.P("")
}
