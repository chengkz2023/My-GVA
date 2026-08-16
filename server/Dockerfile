FROM golang:alpine as builder

WORKDIR /go/src/github.com/chengkz2023/My-GVA/server
COPY . .

RUN go env -w GO111MODULE=on \
    && go env -w GOPROXY=https://goproxy.cn,direct \
    && go env -w CGO_ENABLED=0 \
    && go env \
    && go mod tidy \
    && go build -o server .

FROM alpine:latest

LABEL MAINTAINER="boyking-admin"
# 设置时区
ENV TZ=Asia/Shanghai
RUN apk update && apk add --no-cache tzdata openntpd \
    && ln -sf /usr/share/zoneinfo/$TZ /etc/localtime && echo $TZ > /etc/timezone

WORKDIR /go/src/github.com/chengkz2023/My-GVA/server

COPY --from=0 /go/src/github.com/chengkz2023/My-GVA/server/server ./
COPY --from=0 /go/src/github.com/chengkz2023/My-GVA/server/config.yaml ./

# 运行时挂载配置与日志目录，例如：
# docker run -d -v /宿主机/config.yaml:/go/src/github.com/chengkz2023/My-GVA/server/config.yaml \
#   -v /宿主机/log:/go/src/github.com/chengkz2023/My-GVA/server/log -p 8888:8888 boyking-admin/server:latest

EXPOSE 8888
ENTRYPOINT ./server -c config.yaml
