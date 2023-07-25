FROM golang:1.18-alpine AS builder
RUN sed -i 's/https:\/\/dl-cdn.alpinelinux.org/http:\/\/nexus.clouditera.com:8081\/repository\/apk/g' /etc/apk/repositories && \
    apk update && apk add tzdata dmidecode
RUN cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime
RUN echo Asia/Shanghai > /etc/timezone
RUN apk add git make

ENV CGO_ENABLED 0
ENV GOPROXY=http://nexus.clouditera.com:8081/repository/goproxy/,direct
ENV APP_NAME=message

WORKDIR /build
ADD go.mod .
ADD go.sum .
RUN go mod download
COPY . .
ARG MESSAGE_BRANCH MESSAGE_COMMIT_ID
ENV BRANCH=${MESSAGE_BRANCH}
ENV COMMIT_ID=${MESSAGE_COMMIT_ID}
RUN make build


FROM alpine:3.18
RUN sed -i 's/https:\/\/dl-cdn.alpinelinux.org/http:\/\/nexus.clouditera.com:8081\/repository\/apk/g' /etc/apk/repositories && \
    apk update && apk add tzdata dmidecode
RUN cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime
RUN echo Asia/Shanghai > /etc/timezone
RUN apk del tzdata
WORKDIR /prod/bin
RUN apk update && apk add supervisor

COPY --from=builder /build/bin/message /prod/bin/message
COPY --from=builder /build/script/wrapper.sh /prod/bin/wrapper.sh
COPY --from=builder /build/supervisorconf/* /etc/supervisor/conf.d/
COPY --from=builder /build/supervisorconf/supervisord.conf  /etc/supervisord.conf

CMD ["/usr/bin/supervisord", "-n", "-c", "/etc/supervisord.conf"]
