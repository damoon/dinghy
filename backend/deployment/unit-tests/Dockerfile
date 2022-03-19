FROM golang:1.18.0-buster@sha256:a22f80e34a9f21a7e11a8ef8d29e04366da6993bf4f8efa41202899250130e6c
WORKDIR /app

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY cmd cmd
COPY pkg pkg

CMD ( go test ./... && echo "done" ) || echo "failed"
