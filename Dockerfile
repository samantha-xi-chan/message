FROM rd.clouditera.com/infra/golang_bizcache:1.18.10 AS builder
# FROM golang:1.18.10 AS builder

WORKDIR /app
COPY . /app/src
RUN cd /app/src && \
    go env -w GO111MODULE=on && \
    go env && \
    go env -w GOPROXY=https://goproxy.cn,direct && \
    go mod tidy && \
    sh script/build.sh 

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/main /app/main
EXPOSE 10051 9102 18081
CMD ["./main"]

