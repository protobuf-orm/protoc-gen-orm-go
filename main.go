package main

import (
	"flag"

	"google.golang.org/protobuf/compiler/protogen"
)

func main() {
	h := Handler{}

	var flags flag.FlagSet
	flags.StringVar(&h.Store.Name, "store.name", "store.g.go", "output filename of server and client definitions")
	flags.StringVar(&h.Query.Namer, "query.namer", "{{ .Name }}.g.go", "golang text template for output filename of query utils")

	opts := protogen.Options{ParamFunc: flags.Set}
	opts.Run(h.Run)
}
