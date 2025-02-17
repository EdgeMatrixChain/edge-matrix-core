PROJECTNAME=$(shell basename "$(PWD)")
VERSION=-ldflags="-X main.Version=$(shell git describe --tags)"

.PHONY: help run build install license
all: help

help: Makefile
	@echo
	@echo "Choose a make command to run in "$(PROJECTNAME)":"
	@$(MAKE) -pRrq -f $(lastword $(MAKEFILE_LIST)) : 2>/dev/null | awk -v RS= -F: '/^# File/,/^# Finished Make data base/ {if ($$1 !~ "^[#.]") {print $$1}}' | sort | egrep -v -e '^[^[:alnum:]]' -e '^$@$$'
	@echo

check-git:
	@which git > /dev/null || (echo "git is not installed. Please install and try again."; exit 1)

check-go:
	@which go > /dev/null || (echo "Go is not installed.. Please install and try again."; exit 1)

check-protoc:
	@which protoc > /dev/null || (echo "protoc is not installed. Please install and try again."; exit 1)

get:
	@echo "  >  \033[32mDownloading & Installing all the modules...\033[0m "
	go mod tidy && go mod download

protoc:
	protoc --go_out=./core --go-grpc_out=./core ./core/application/proto/*.proto
	protoc --go_out=./core --go-grpc_out=./core ./core/network/proto/*.proto

