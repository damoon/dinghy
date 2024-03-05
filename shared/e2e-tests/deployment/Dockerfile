FROM golang:1.22.1-bookworm@sha256:6699d2852712f090399ccd4e8dfd079b4d55376f3ab3aff5b2dc8b7b1c11e27e AS build-env

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
