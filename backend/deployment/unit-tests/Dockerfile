FROM golang:1.17.8-buster@sha256:0b49d92e9b33cc1972a84942e1a4504038afa73290f73d4b95e00d10e47999f3
WORKDIR /app

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY cmd cmd
COPY pkg pkg

CMD ( go test ./... && echo "done" ) || echo "failed"
