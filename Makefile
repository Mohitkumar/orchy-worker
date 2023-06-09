SHELL := /bin/bash
BINARY_NAME=orchy
GOPATH ?= $(shell go env GOPATH)

GO                  := GO111MODULE=on go
compile:
	protoc api/v1/*.proto \
		--go_out=. \
		--go-grpc_out=. \
		--go_opt=paths=source_relative \
		--go-grpc_opt=paths=source_relative \
		--proto_path=.