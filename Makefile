build:
	go build -o bin/server cmd/server/main.go
	go build -o bin/client cmd/client/main.go

run-server:
	go run cmd/server/main.go

run-client:
	go run cmd/client/main.go

generate-proto:
	protoc --go_out=./internal --plugin=protoc-gen-go=$(shell go tool -n protoc-gen-go) --go_opt=paths=source_relative --proto_path=./internal ./internal/message.proto