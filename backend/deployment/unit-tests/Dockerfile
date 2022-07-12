FROM golang:1.18.3-buster@sha256:7acdbec4276ce3c8915ffa03148bd4b41dfb46a190fd1242e9df5b11061d282d
WORKDIR /app

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY cmd cmd
COPY pkg pkg

CMD ( go test ./... && echo "done" ) || echo "failed"
