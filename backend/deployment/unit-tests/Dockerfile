FROM golang:1.18.1-buster@sha256:a58a53ecc3088e9ab801e0a8851dae8294f9931446685c43dda9646144f38af3
WORKDIR /app

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY cmd cmd
COPY pkg pkg

CMD ( go test ./... && echo "done" ) || echo "failed"
