FROM golang:1.17.5-buster@sha256:47030a15ca7cccf40da9f968fabebc5e99ce5455ad4a688cda3d748278a9c967
WORKDIR /app

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY cmd cmd
COPY pkg pkg

CMD ( go test ./... && echo "done" ) || echo "failed"
