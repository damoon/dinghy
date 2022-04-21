FROM golang:1.18.1-buster@sha256:6ed1375f000a9086969581dd1dcc676495599228fa6b27c3219c4a76b8f09cc4
WORKDIR /app

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY cmd cmd
COPY pkg pkg

CMD ( go test ./... && echo "done" ) || echo "failed"
