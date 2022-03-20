FROM golang:1.18.0-buster@sha256:35cf8e5001bd5b2c3fc0286e379bed24cc0dd1b5b5d50ddd524b30e09696b0cc
WORKDIR /app

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY cmd cmd
COPY pkg pkg

CMD ( go test ./... && echo "done" ) || echo "failed"
