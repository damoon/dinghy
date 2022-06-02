FROM golang:1.18.3-buster@sha256:8ec5e989db0ae2aaeb0015cb8cb0a868d262d188c15bca06f61436c6f2ce5559
WORKDIR /app

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY cmd cmd
COPY pkg pkg

CMD ( go test ./... && echo "done" ) || echo "failed"
