FROM golang:1.18.5-buster@sha256:cb94fa4e6f97e8ce76a1809e4e1d071a39753c1bc696035a1992b9dc3d444ebc
WORKDIR /app

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY cmd cmd
COPY pkg pkg

CMD ( go test ./... && echo "done" ) || echo "failed"
