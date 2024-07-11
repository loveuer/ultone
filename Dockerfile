#FROM golang:1.20-alpine AS builder
FROM repository.umisen.com/external/golang:latest AS builder

ENV GO111MODULE on
ENV CGO_ENABLED 0
ENV GOOS linux
ENV GOPROXY https://goproxy.io

WORKDIR /build

COPY . .

RUN go mod download
RUN go build -ldflags '-s -w' -o server .

#FROM alpine:latest
FROM repository.umisen.com/external/alpine:latest

RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.tuna.tsinghua.edu.cn/g' /etc/apk/repositories && apk add curl

ENV TZ Asia/Shanghai

WORKDIR /app

RUN mkdir -p /data

COPY --from=builder /build/server /app/server
COPY etc /app/etc

CMD ["/app/server", "-c", "/app/etc/config.json"]
