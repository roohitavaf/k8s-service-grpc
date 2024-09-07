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
	rm -f bin/echo-client
	rm -f bin/echo-server

deploy-headless:
	kubectl apply -f config/headless-service.yaml
	ko apply -f config/echo-server.yaml
	ko apply -f config/echo-client-headless.yaml


deploy-cluster-ip:
	kubectl apply -f config/cluster-ip-service.yaml
	ko apply -f config/echo-server.yaml
	ko apply -f config/echo-client-cluster-ip.yaml

undeploy-all:
	kubectl delete -f config/

.PHONY: proto build clean deploy-headless deploy-cluster-ip undeploy-all