FROM golang:1.21.3-bookworm@sha256:5bafbbb109f02aaf6b41ddc19f54919773c3006f1cbda1599112603367642f0e
WORKDIR /app

COPY go.mod .
COPY go.sum .
RUN GOPROXY=https://goproxy.io go mod download

COPY cmd cmd
COPY pkg pkg

CMD ( GOPROXY=https://goproxy.io go test ./... && echo "done" ) || echo "failed"
