FROM        golang:1.12-alpine3.10 AS build
MAINTAINER  Shaddock <hushuang@gmail.com>
WORKDIR     /go/src/github.com/ntfs32/wework-bot


RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories && \
    apk add --update -t build-deps curl libc-dev gcc libgcc git wget && \
    wget -q https://github.com/upx/upx/releases/download/v3.95/upx-3.95-amd64_linux.tar.xz && \
    tar xvJf upx-3.95-amd64_linux.tar.xz && \
    cp upx-3.95-amd64_linux/upx /usr/local/bin/ && \
    chmod +x /usr/local/bin/upx

ENV GO111MODULE=on
ENV GOPROXY=https://goproxy.io
ADD . .
RUN for obj in `ls lib/`; do echo "build ./lib/$obj/*.go into /opt/build/$obj ..."; go build -o /opt/build/$obj ./lib/$obj/*.go; upx /opt/build/$obj; done && \
    chmod +x /opt/build/* && \
    apk del --purge build-deps && \
    rm -rf /var/cache/apk/* && \
    rm -rf /go



FROM almir/webhook

RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories && \
    apk add bash vim

WORKDIR /opt/
ADD hooks /opt/hooks
ADD scripts /opt/scripts

COPY  --from=build /opt/build/* /opt/scripts/
RUN ls /opt/hooks && chmod +x /opt/scripts/*

EXPOSE 9000
ENTRYPOINT ["webhook", "-hooks", "/opt/hooks/hooks.yaml", "-verbose"]