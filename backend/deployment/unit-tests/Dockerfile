FROM golang:1.20.2-buster@sha256:5f9d227a93a4453df44827dbf44acef84e0b0c3abf39470d8a69b28bbd56733f
WORKDIR /app

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY cmd cmd
COPY pkg pkg

CMD ( go test ./... && echo "done" ) || echo "failed"
