FROM golang:1.14.1 AS workbench
ENV ENVIRONMENT development
RUN mkdir /tests
WORKDIR /tests

RUN apt-get update
RUN apt-get install -y bats curl

ADD https://dl.google.com/linux/direct/google-chrome-stable_current_amd64.deb /tests/chrome.deb
RUN apt install -y ./chrome.deb

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY main.go /tests/main.go
RUN go build -o /go/bin/main .

COPY browser.bats /tests
COPY favicon.png /tests

ENTRYPOINT [ "bats", "." ]
