FROM golang:1.17.3-buster@sha256:781755154df17bf0b91f02ed0381ab31b46b6a9ca7d3c93f2a55e491af19970c
WORKDIR /app

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY cmd cmd
COPY pkg pkg

CMD ( go test ./... && echo "done" ) || echo "failed"
