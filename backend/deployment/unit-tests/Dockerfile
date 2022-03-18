FROM golang:1.18.0-buster@sha256:4d2553850f501c598d05e9a587653a5a81d6c7967f4368ade8e492fcf250f8f1
WORKDIR /app

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY cmd cmd
COPY pkg pkg

CMD ( go test ./... && echo "done" ) || echo "failed"
