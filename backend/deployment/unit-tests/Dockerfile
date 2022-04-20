FROM golang:1.18.1-buster@sha256:279e4ec837160ba80a6ce81e8ae5a5fa4723dc720e61e9b262b82df2c9b53b0f
WORKDIR /app

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY cmd cmd
COPY pkg pkg

CMD ( go test ./... && echo "done" ) || echo "failed"
