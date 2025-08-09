package app

func (w *Work) xClient() {
	w.P("type Client interface {")
	for _, v := range w.Entities {
		w.P("	", v.Name(), "() ", w.Package.Ident(v.Name()+"ServiceClient"))
	}
	w.P("}")
	w.P("")

	w.P("func NewClient(c *", grpc.Ident("ClientConn"), ") Client {")
	w.P("	return &client{")
	for _, v := range w.Entities {
		w.P("	_", v.Name(), ": ", w.Package.Ident("New"+v.Name()+"ServiceClient"), "(c),")
	}
	w.P("	}")
	w.P("}")
	w.P("")

	w.P("type client struct {")
	for _, v := range w.Entities {
		w.P("	_", v.Name(), " ", w.Package.Ident(v.Name()+"ServiceClient"))
	}
	w.P("}")
	w.P("")
	for _, v := range w.Entities {
		w.P("func (c *client) ", v.Name(), "() ", w.Package.Ident(v.Name()+"ServiceClient"), " { return c._", v.Name(), " }")
	}
	w.P("")
}
