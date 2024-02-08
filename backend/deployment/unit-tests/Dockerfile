FROM golang:1.22.0-bookworm@sha256:874c2677e43be36a429823f2742af85844772664f273c1c8c8235f10aba51cd3
WORKDIR /app

COPY go.mod .
COPY go.sum .
RUN GOPROXY=https://goproxy.io go mod download

COPY cmd cmd
COPY pkg pkg

CMD ( GOPROXY=https://goproxy.io go test ./... && echo "done" ) || echo "failed"
