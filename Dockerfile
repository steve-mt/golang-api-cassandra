FROM golang:1.8-alpine

MAINTAINER Steve Azzopardi <steveazz@outlook.com>

COPY . /go/src/github.com/SteveAzz/stream-api

WORKDIR /go/src/github.com/SteveAzz/stream-api

RUN go install -v

ENTRYPOINT /go/bin/stream-api

EXPOSE 8080