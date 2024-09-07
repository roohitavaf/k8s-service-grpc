KO_DOCKER_REPO ?= kind.local

export KO_DOCKER_REPO
export KO_DEFAULTBASEIMAGE=alpine:latest

proto:
	protoc --proto_path=api/v1 --go_out=pkg/echo --go_opt=paths=source_relative \
	--go-grpc_out=pkg/echo --go-grpc_opt=paths=source_relative \
	api/v1/echo.proto

build: proto
	go build -o bin/echo-client cmd/echo-client/*
	go build -o bin/echo-server cmd/echo-server/*

clean:
	rm -f pkg/echo/*.pb.go
	rm -f bin/myapp

publish:
	ko publish -t latest ./cmd/echo-server
	ko publish -t latest ./cmd/echo-client

server:
	ko apply -f config/echo-server.yaml

client:
	ko apply -f config/echo-client.yaml

.PHONY: proto build clean