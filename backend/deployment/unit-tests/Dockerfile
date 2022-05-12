FROM golang:1.18.2-buster@sha256:949a195c5d44e7b7f71499c146def008c7e98b8d02c276e567ab2e5c6971587b
WORKDIR /app

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY cmd cmd
COPY pkg pkg

CMD ( go test ./... && echo "done" ) || echo "failed"
