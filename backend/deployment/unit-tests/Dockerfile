FROM golang:1.20.1-buster@sha256:9628a1ad396ad02e6d26b8506cd36b2e04c56f827f060bd06fd095963263fea2
WORKDIR /app

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY cmd cmd
COPY pkg pkg

CMD ( go test ./... && echo "done" ) || echo "failed"
