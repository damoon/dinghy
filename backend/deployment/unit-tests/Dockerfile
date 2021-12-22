FROM golang:1.17.5-buster@sha256:ee9771f71e7f70eb8c09852ab05844d9e896a00f931af40fbb5b9481c732230c
WORKDIR /app

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY cmd cmd
COPY pkg pkg

CMD ( go test ./... && echo "done" ) || echo "failed"
