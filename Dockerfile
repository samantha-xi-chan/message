FROM golang:1.18-alpine AS builder
RUN apk add tzdata dmidecode
RUN cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime
RUN echo Asia/Shanghai > /etc/timezone
RUN apk add git make

ENV CGO_ENABLED 0
ENV GOPROXY https://goproxy.cn,direct
ENV APP_NAME=message
ENV MONGO_URL=mongodb://mongo:27017
ENV AMQP_URL=amqp://guest:guest@rabbitmq:5672/

WORKDIR /build
ADD go.mod .
ADD go.sum .
RUN go mod download
COPY . .
RUN make build


FROM alpine
RUN apk add tzdata dmidecode
RUN cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime
RUN echo Asia/Shanghai > /etc/timezone
RUN apk del tzdata
WORKDIR /prod/bin
RUN apk update && apk add supervisor

COPY --from=builder /build/bin/message_linux_amd64 /prod/bin/message
COPY --from=builder /build/supervisorconf/*  /etc/supervisor/conf.d/
COPY --from=builder /build/supervisorconf/supervisord.conf  /etc/supervisord.conf

CMD ["/usr/bin/supervisord", "-n", "-c", "/etc/supervisord.conf"]
