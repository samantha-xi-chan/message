BUILD_DATE:=$(shell date '+%Y%m%d%H%M%S')
BRNACH:=$(shell git rev-parse --abbrev-ref HEAD)
COMMIT_ID:=$(shell git rev-parse HEAD)
OUTPUT_DIR:=./bin

LDFLAGS:=-s -X main.BUILD_DATE=${BUILD_DATE} -X main.GIT_BRANCH=${BRNACH} -X main.GIT_COMMIT=${COMMIT_ID}

ifneq (${MONGO_URL},)
	LDFLAGS:=${LDFLAGS} -X internal.MONGO_URL=${MONGO_URL}
endif

ifneq (${AMQP_URL},)
	LDFLAGS:=${LDFLAGS} -X internal.AMQP_URL=${AMQP_URL}
endif

.PHONY: all build run gotool clean help

all: gotool build

help:
	@echo "make help/all"

pb:
	protoc api/proto/*.proto --go_out=plugins=grpc:.

build:
	@if [ -d bin ] ; then rm -rf ./bin ; fi
	mkdir bin

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
