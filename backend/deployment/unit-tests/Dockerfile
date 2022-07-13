FROM golang:1.18.4-buster@sha256:295d17b21f1b206054dcd72fa1a1b26f197fc661b0901a312ee471e40b251757
WORKDIR /app

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY cmd cmd
COPY pkg pkg

CMD ( go test ./... && echo "done" ) || echo "failed"
