FROM golang:1.8.3-alpine3.6

RUN apk update && apk add \
  gcc \
  git \
  linux-headers \
  make \
  musl-dev

RUN go get -u github.com/golang/lint/golint

ADD . /go/src/github.com/simonferquel/integration-kit/
WORKDIR /go/src/github.com/simonferquel/integration-kit
