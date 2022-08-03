FROM golang:1.18.5-buster@sha256:eb4b3bd59da4604b0ab27d3837e39dacfba4a48cdd289419a3c49993c0eadba4
WORKDIR /app

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY cmd cmd
COPY pkg pkg

CMD ( go test ./... && echo "done" ) || echo "failed"
