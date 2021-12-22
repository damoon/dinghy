FROM golang:1.17.5-buster@sha256:7f89be6ec02f7eb4b429e9a43f4ac7d49d2c0cb79a99a28e2d1efe1e7e8a3473
WORKDIR /app

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY cmd cmd
COPY pkg pkg

CMD ( go test ./... && echo "done" ) || echo "failed"
