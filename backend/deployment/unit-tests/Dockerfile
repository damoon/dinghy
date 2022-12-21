FROM golang:1.19.4-buster@sha256:c1e438a538be3fe72468f1569a57565225fcab11762f372f46b20d65bba45ab9
WORKDIR /app

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY cmd cmd
COPY pkg pkg

CMD ( go test ./... && echo "done" ) || echo "failed"
