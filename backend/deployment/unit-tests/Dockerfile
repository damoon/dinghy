FROM golang:1.18.0-buster@sha256:e4950365f46f107ba92f8a9557defccb38f151a79fd5453e9424ec34de08c35d
WORKDIR /app

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY cmd cmd
COPY pkg pkg

CMD ( go test ./... && echo "done" ) || echo "failed"
