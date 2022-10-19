.PHONY: all build run gotool clean help

all: gotool build

help:
	@echo "make help/all"

pb:
	protoc api/proto/*.proto --go_out=plugins=grpc:.

build:
	echo "protoc ing"
	@if [ -d bin ] ; then rm -rf ./bin ; fi
	mkdir bin
	echo "Compiling for every OS and Platform"
	CGO_ENABLED=0 GOOS=linux  GOARCH=amd64  go build -o bin/message_linux_amd64  .
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64  go build -o bin/message_darwin_x86   .
	CGO_ENABLED=0 GOOS=darwin GOARCH=arm64  go build -o bin/message_darwin_arm64 .

run:
	go run .

gotool:
	go fmt ./
	go vet ./

clean:
	@if [ -d bin ] ; then rm ./bin/* ; fi

all: gotool build
