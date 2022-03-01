FROM golang:1.17.7-buster@sha256:428a77349bc4b83a13ad950160288e323a6552680c11a7473f776cbd97beae6b AS build-env

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
