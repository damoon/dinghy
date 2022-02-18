FROM golang:1.17.7-buster@sha256:270086796f4fea7905b19889d33a6063e93eb6ad9c6fde11aaf3b2f82a757bef
WORKDIR /app

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY cmd cmd
COPY pkg pkg

CMD ( go test ./... && echo "done" ) || echo "failed"
