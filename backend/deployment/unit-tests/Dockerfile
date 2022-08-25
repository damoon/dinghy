FROM golang:1.19.0-buster@sha256:d84495e2ad12dfeb51dec3c9f170659ebd09db234d6990b3364f877785684a14
WORKDIR /app

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY cmd cmd
COPY pkg pkg

CMD ( go test ./... && echo "done" ) || echo "failed"
