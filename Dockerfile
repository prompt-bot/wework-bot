FROM        golang:1.12-alpine3.10 AS build
MAINTAINER  Shaddock <hushuang@gmail.com>
WORKDIR     /go/src/github.com/ntfs32/wework-bot


RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories && \
    apk add --update -t build-deps curl libc-dev gcc libgcc git
ENV GO111MODULE=on
ADD . .
RUN go build  -o /opt/build/gitlab ./gitlab/gitlab.go && \
    apk del --purge build-deps && \
    rm -rf /var/cache/apk/* && \
    rm -rf /go



FROM almir/webhook

RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories && \
    apk add bash vim

WORKDIR /opt/
ADD hooks /opt/hooks
ADD scripts /opt/scripts

COPY  --from=build /opt/build/* /usr/local/bin/
RUN ls /opt/hooks/ && chmod +x /opt/scripts/*

EXPOSE 9000
ENTRYPOINT ["webhook", "-hooks", "/opt/hooks/hooks.yaml", "-verbose"]