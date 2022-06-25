FROM golang:1.18.3-buster@sha256:3ad5048c0c63cbc32d1c576bda513ebfd6e1914b6cd57837e0cd0a31fe81de69
WORKDIR /app

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY cmd cmd
COPY pkg pkg

CMD ( go test ./... && echo "done" ) || echo "failed"
