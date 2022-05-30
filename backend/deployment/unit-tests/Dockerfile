FROM golang:1.18.2-buster@sha256:ad5d8648c50bd3b2dfdb3bd4f6ab52a82b43d949f53c1998884b1653e0b873b5
WORKDIR /app

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY cmd cmd
COPY pkg pkg

CMD ( go test ./... && echo "done" ) || echo "failed"
