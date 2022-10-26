FROM golang:1.19.2-buster@sha256:202b7098b6bd5c91faacce750514a61e2dd2c31a93b419eca3bb0af3c409835b
WORKDIR /app

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY cmd cmd
COPY pkg pkg

CMD ( go test ./... && echo "done" ) || echo "failed"
