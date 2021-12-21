FROM golang:1.17.5-buster@sha256:8f73326871cfb1482549190da48e4e29e06939e2ba6830c8cf7fc551d252f7cb
WORKDIR /app

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY cmd cmd
COPY pkg pkg

CMD ( go test ./... && echo "done" ) || echo "failed"
