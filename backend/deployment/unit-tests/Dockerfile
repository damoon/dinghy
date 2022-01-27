FROM golang:1.17.6-buster@sha256:d3ca64c1e652d7beee33930b0aabb2ab668097690b66698208551e700dbb2e3a
WORKDIR /app

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY cmd cmd
COPY pkg pkg

CMD ( go test ./... && echo "done" ) || echo "failed"
