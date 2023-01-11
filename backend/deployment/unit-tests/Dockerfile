FROM golang:1.19.5-buster@sha256:33d767d032a989ba380e6448b512e9b7056a29a3efbf0409427fdd8445f392e5
WORKDIR /app

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY cmd cmd
COPY pkg pkg

CMD ( go test ./... && echo "done" ) || echo "failed"
