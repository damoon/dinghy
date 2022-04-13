FROM golang:1.18.1-buster@sha256:2b52df3e45d4f2c8ac8d6db12366d09394189bc13d67c27eb2b44b79dfa4db00
WORKDIR /app

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY cmd cmd
COPY pkg pkg

CMD ( go test ./... && echo "done" ) || echo "failed"
