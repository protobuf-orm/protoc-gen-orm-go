module github.com/protobuf-orm/protoc-gen-orm-go

go 1.27.0

require (
	github.com/lesomnus/protobuf-patch v0.0.0-20260803175157-e1b7a0c2804f
	github.com/protobuf-orm/protobuf-orm v0.0.0-20260901155106-6bf45a2a1e67
	google.golang.org/grpc v1.74.2
	google.golang.org/protobuf v1.36.11
)

require (
	github.com/ettle/strcase v0.2.0 // indirect
	github.com/protobuf-orm/protoc-gen-orm-service v0.0.0-20260901155246-12064a5a7fa8 // indirect
	golang.org/x/net v0.40.0 // indirect
	golang.org/x/sys v0.33.0 // indirect
	golang.org/x/text v0.25.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20250528174236-200df99c418a // indirect
	google.golang.org/grpc/cmd/protoc-gen-go-grpc v1.5.1 // indirect
)

tool (
	github.com/protobuf-orm/protoc-gen-orm-service
	google.golang.org/grpc/cmd/protoc-gen-go-grpc
	google.golang.org/protobuf/cmd/protoc-gen-go
)
