module github.com/BrunoCupertino/stock-hollywood

go 1.24.2

// google.golang.org/grpc/cmd/protoc-gen-go-grpc
tool google.golang.org/protobuf/cmd/protoc-gen-go

// google.golang.org/grpc/cmd/protoc-gen-go-grpc v1.5.1 // indirect
require google.golang.org/protobuf v1.36.10

require (
	github.com/anthdm/hollywood v1.0.5
	github.com/google/uuid v1.6.0
)

require (
	github.com/DataDog/gostackparse v0.7.0 // indirect
	github.com/klauspost/cpuid/v2 v2.0.9 // indirect
	github.com/zeebo/xxh3 v1.0.2 // indirect
)
