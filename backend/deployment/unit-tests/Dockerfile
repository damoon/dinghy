FROM golang:1.19.2-buster@sha256:d7067103aa14d258d61c6ac1c4903779653fb43211346f267f2bf2e86d0212fc
WORKDIR /app

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY cmd cmd
COPY pkg pkg

CMD ( go test ./... && echo "done" ) || echo "failed"
