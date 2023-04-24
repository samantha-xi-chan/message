BUILD_DATE:=$(shell date '+%Y%m%d%H%M%S')
ifeq ($(strip $(BRANCH)),)
    BRANCH := $(shell git rev-parse --abbrev-ref HEAD)
else
	BRANCH := $(BRANCH)
endif
ifeq ($(strip $(COMMIT_ID)),)
    COMMIT_ID := $(shell git rev-parse --abbrev-ref HEAD)
else
	COMMIT_ID := $(COMMIT_ID)
endif
OUTPUT_DIR:=./bin

LDFLAGS:=-s -X main.BUILD_DATE=${BUILD_DATE} -X main.GIT_BRANCH=${BRANCH} -X main.GIT_COMMIT=${COMMIT_ID}

.PHONY: all build run gotool clean help

all: gotool build

help:
	@echo "make help/all"

pb:
	protoc api/proto/*.proto --go_out=plugins=grpc:.

build:
	CGO_ENABLED=0 GOOS=linux  GOARCH=amd64 go build -buildvcs=false -v -ldflags "${LDFLAGS}" -o ${OUTPUT_DIR}/message

all_platform:
	@echo "Compiling for every OS and Platform"
	CGO_ENABLED=0 GOOS=linux  GOARCH=amd64 go build -v -ldflags "${LDFLAGS}" -o ${OUTPUT_DIR}/message_linux_amd64
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -v -ldflags "${LDFLAGS}" -o ${OUTPUT_DIR}//message_darwin_x86
	CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 go build -v -ldflags "${LDFLAGS}" -o ${OUTPUT_DIR}//message_darwin_arm64

run:
	go run .

gotool:
	go fmt ./
	go vet ./

clean:
	@if [ -d bin ] ; then rm ./bin/* ; fi

all: gotool build
