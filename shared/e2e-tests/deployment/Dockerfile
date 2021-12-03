FROM golang:1.17.3-buster@sha256:781755154df17bf0b91f02ed0381ab31b46b6a9ca7d3c93f2a55e491af19970c AS build-env

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
