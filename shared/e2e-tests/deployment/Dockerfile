FROM golang:1.21.5-bookworm@sha256:24e832041c0df9265032b46ce2eb4a5e887119089b84705565f6cf4bd534caf6 AS build-env

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
