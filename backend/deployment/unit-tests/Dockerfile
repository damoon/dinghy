FROM golang:1.17.3-buster@sha256:6d0772539cd83e5a8b24a0003e82e042148f89762d36c610c220536322541e94
WORKDIR /app

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY cmd cmd
COPY pkg pkg

CMD ( go test ./... && echo "done" ) || echo "failed"
