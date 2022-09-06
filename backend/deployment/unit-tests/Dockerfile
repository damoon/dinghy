FROM golang:1.19.1-buster@sha256:d71125b5d1e299a8d9829eb09dbeb0cf576c461cbded210a9a6242d64072adc9
WORKDIR /app

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY cmd cmd
COPY pkg pkg

CMD ( go test ./... && echo "done" ) || echo "failed"
