FROM golang:1.17.2-buster@sha256:9cf49bf1e8cf0c8d995899546bea5d876e59b5a9fdeb1575cb6197f97973c20e
WORKDIR /app

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY cmd cmd
COPY pkg pkg

CMD ( go test ./... && echo "done" ) || echo "failed"
