FROM golang:1.19.1-buster@sha256:7656e693167ae561f883fe292ad3e27e05eb681bbfb53ebe0fbb31c95cbc4967
WORKDIR /app

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY cmd cmd
COPY pkg pkg

CMD ( go test ./... && echo "done" ) || echo "failed"
