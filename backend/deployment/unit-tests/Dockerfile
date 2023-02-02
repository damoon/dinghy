FROM golang:1.20.0-buster@sha256:fc1e0b02674cef154ceca0ba6ca4d3658bd385d5faadecbbd44809727a56329d
WORKDIR /app

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY cmd cmd
COPY pkg pkg

CMD ( go test ./... && echo "done" ) || echo "failed"
