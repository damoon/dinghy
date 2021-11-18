FROM golang:1.17.3-buster@sha256:099de80fda2a712d3ee9f66f7f8b62061aad4b5e0ff43f66c86f16ba6a3a89cd
WORKDIR /app

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY cmd cmd
COPY pkg pkg

CMD ( go test ./... && echo "done" ) || echo "failed"
