FROM golang:1.18.1-buster@sha256:d8f864bac466e488bb425a59ff9bc104828cd44d0985a8d2fcb5570a9250258a
WORKDIR /app

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY cmd cmd
COPY pkg pkg

CMD ( go test ./... && echo "done" ) || echo "failed"
