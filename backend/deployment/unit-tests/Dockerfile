FROM golang:1.17.7-buster@sha256:3d11eb4c0032307bbecbb40c30c0719006f4960f8ae4d01383e6b442425a9057
WORKDIR /app

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY cmd cmd
COPY pkg pkg

CMD ( go test ./... && echo "done" ) || echo "failed"
