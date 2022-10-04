FROM golang:1.19.2-buster@sha256:535a9d6a746bfdb02e6c2b506c62aaeb72b8327bcc7e17a5f2cc831a6b4026c2
WORKDIR /app

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY cmd cmd
COPY pkg pkg

CMD ( go test ./... && echo "done" ) || echo "failed"
