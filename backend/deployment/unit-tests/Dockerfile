FROM golang:1.17.6-buster@sha256:a5ba2c081bc7174317ce7f54cf5689c2365454b11a5e271a911d6a2edec1163d
WORKDIR /app

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY cmd cmd
COPY pkg pkg

CMD ( go test ./... && echo "done" ) || echo "failed"
