FROM golang:1.18.2-buster@sha256:e5013e70182658129a37b368718be0d12deff03655a1514854ac9b1d54d8d94c
WORKDIR /app

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY cmd cmd
COPY pkg pkg

CMD ( go test ./... && echo "done" ) || echo "failed"
