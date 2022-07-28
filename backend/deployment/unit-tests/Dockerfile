FROM golang:1.18.4-buster@sha256:6960d62610b18b7224d2c5572b4bb177890b9ab7bf70ebaf34e2e9ca662a46e9
WORKDIR /app

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY cmd cmd
COPY pkg pkg

CMD ( go test ./... && echo "done" ) || echo "failed"
