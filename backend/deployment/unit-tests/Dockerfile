FROM golang:1.20.2-buster@sha256:b089a60efbf123ecfa2b0d8b955c9f441d83d08918dd358a7a8713e71c17c1d7
WORKDIR /app

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY cmd cmd
COPY pkg pkg

CMD ( go test ./... && echo "done" ) || echo "failed"
