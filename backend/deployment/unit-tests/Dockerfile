FROM golang:1.18.0-buster@sha256:284da7e6ad7fd1250a785c24912fb6f8652176b588d637de5c5b3d4a4aa1693a
WORKDIR /app

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY cmd cmd
COPY pkg pkg

CMD ( go test ./... && echo "done" ) || echo "failed"
