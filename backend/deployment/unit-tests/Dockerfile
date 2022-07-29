FROM golang:1.18.4-buster@sha256:477b10a289d7be9c5124da52f2ff86d13c24e591b82da3b6b1ad3b8f1a04c8d9
WORKDIR /app

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY cmd cmd
COPY pkg pkg

CMD ( go test ./... && echo "done" ) || echo "failed"
