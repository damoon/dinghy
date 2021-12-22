FROM golang:1.17.5-buster@sha256:62a929434e43b91ba7a215c1233d97d0649f0d6c3dc2cb03d0c24ac0f81f1982
WORKDIR /app

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY cmd cmd
COPY pkg pkg

CMD ( go test ./... && echo "done" ) || echo "failed"
