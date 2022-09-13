FROM golang:1.19.1-buster@sha256:7656e693167ae561f883fe292ad3e27e05eb681bbfb53ebe0fbb31c95cbc4967 AS build-env

RUN mkdir /tests
WORKDIR /tests

RUN apt-get update
RUN apt-get install -y bats curl

COPY deployment/run-tests.sh /usr/local/bin/run-tests.sh
CMD run-tests.sh

ADD https://dl.google.com/linux/direct/google-chrome-stable_current_amd64.deb /tests/chrome.deb
RUN apt install -y ./chrome.deb

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY main.go /tests/main.go
RUN go build -o /go/bin/main .

COPY . /tests
