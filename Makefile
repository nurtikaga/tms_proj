DATE?=$(shell date -u "+%Y-%m-%d %H:%M:%S")
APP_NAME?=tms-project
APP_VERSION?=latest
LDFLAGS=-s -w -X 'main.AppName=${APP_NAME}' -X 'main.AppVersion=${APP_VERSION}' -X 'main.BuildDate=${DATE}'
DOCKER=docker

.PHONY: build
build: 
protoc --go_out=go --go_opt=paths=source_relative --go-grpc_out=go --go-grpc_opt=paths=source_relative ShipmentManaging.proto