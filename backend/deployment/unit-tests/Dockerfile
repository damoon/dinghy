FROM golang:1.22.2-bookworm@sha256:48b942a5a9803fafeb748f1a9c6edeb71e06653bb4df25f63f909151e15e9618
WORKDIR /app

COPY go.mod .
COPY go.sum .
RUN GOPROXY=https://goproxy.io go mod download

COPY cmd cmd
COPY pkg pkg

CMD ( GOPROXY=https://goproxy.io go test ./... && echo "done" ) || echo "failed"
