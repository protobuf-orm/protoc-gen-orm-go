package app

func (w *Work) xServerInterface() {
	w.P("type Server interface {")
	for _, v := range w.Entities {
		w.P("	", v.Name(), "() ", w.Package.Ident(v.Name()+"ServiceServer"))
	}
	w.P("}")
	w.P("")
	w.P("func RegisterServer(",
		/* */ "g *", grpc.Ident("Server"), ", ",
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
