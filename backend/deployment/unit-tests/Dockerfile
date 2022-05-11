FROM golang:1.18.2-buster@sha256:8c7ceffb849c06888d9b8d6ee1314883aa9d9f26ab3aa544e7e0dcac3ec41389
WORKDIR /app

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY cmd cmd
COPY pkg pkg

CMD ( go test ./... && echo "done" ) || echo "failed"
