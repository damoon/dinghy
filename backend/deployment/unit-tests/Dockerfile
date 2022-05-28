FROM golang:1.18.2-buster@sha256:30968aa42e8f6f3cf45d12b22ade8aad14d68183ceec5c310e07c01ce49d7f98
WORKDIR /app

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY cmd cmd
COPY pkg pkg

CMD ( go test ./... && echo "done" ) || echo "failed"
