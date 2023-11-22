FROM golang:1.21.4-bookworm@sha256:211553a2281d0d49ffd7d8c766624bc5428e40840a3238cbdd3c868b1d9819b0 AS build-env

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
RUN GOPROXY=https://goproxy.io go mod download

COPY main.go /tests/main.go
RUN go build -o /go/bin/main .

COPY . /tests
