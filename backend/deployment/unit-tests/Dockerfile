FROM golang:1.17.3-buster@sha256:29f41d09f4e91ac5c516c0b92610dfc372f8713db81feb24d45b67eb9232d260
WORKDIR /app

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY cmd cmd
COPY pkg pkg

CMD ( go test ./... && echo "done" ) || echo "failed"
