package app

func (w *Work) xServerInterface() {
	w.P("type Server interface {")
	for _, v := range w.Entities {
		w.P("	", v.Name(), "() ", w.Package.Ident(v.Name()+"ServiceServer"))
	}
	w.P("}")
	w.P("")
	// A ServiceRegistrar and not a *grpc.Server, which is a strict widening:
	// every *grpc.Server is one, so nothing that called this before has to
	// change. What it buys is the callers that are not a gRPC server at all --
	// a server compiled into a page and speaking a datagram protocol to it, a
	// test that registers into something of its own. The per-service
	// Register<E>ServiceServer that protoc-gen-go-grpc emits already takes the
	// interface; this was the one line that did not.
	w.P("// RegisterServer registers every service of `s` with `g`.")
	w.P("//")
	w.P("// It takes a [grpc.ServiceRegistrar] rather than a *grpc.Server so that a")
	w.P("// server which is not gRPC's own can be handed the same set of services.")
	w.P("func RegisterServer(",
		/* */ "g ", grpc.Ident("ServiceRegistrar"), ", ",
		/* */ "s Server",
		") {")
	for _, v := range w.Entities {
		w.P("	", w.Package.Ident("Register"+v.Name()+"ServiceServer"), "(g, s.", v.Name(), "())")
	}
	w.P("}")
	w.P("")
}

func (w *Work) xUnimplementedServerStruct() {
	w.P("type UnimplementedServer struct {")
	for _, v := range w.Entities {
		w.P("	", v.Name(), "Server ", w.Package.Ident(v.Name()+"ServiceServer"))
	}
	w.P("}")
	w.P("")
	for _, v := range w.Entities {
		w.P("func (UnimplementedServer) ", v.Name(), "() ", w.Package.Ident(v.Name()+"ServiceServer"), "{ return Unimplemented", v.Name(), "ServiceServer{} }")
	}

	w.P("")
}

func (w *Work) xStaticServerStruct() {
	w.P("type StaticServer struct {")
	for _, v := range w.Entities {
		w.P("	", v.Name(), "Server ", w.Package.Ident(v.Name()+"ServiceServer"))
	}
	w.P("}")
	w.P("")
	// A value receiver, like UnimplementedServer's, so that a StaticServer is a
	// Server whether or not it is behind a pointer. [Find] matches on the
	// dynamic type, so a stack holding one of the two cannot be asked for the
	// other, and there is no reason for this one to be the awkward kind.
	for _, v := range w.Entities {
		w.P("func (s StaticServer) ", v.Name(), "() ", w.Package.Ident(v.Name()+"ServiceServer"), "{ return s.", v.Name(), "Server }")
	}

	w.P("")
}
