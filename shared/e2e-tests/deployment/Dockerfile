FROM golang:1.17.3-buster@sha256:48f6285a02b7682373e97f153238f3a831268ca2e7dd6204bf540b6801ac63bc AS build-env

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
