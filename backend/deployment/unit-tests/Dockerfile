FROM golang:1.21.3-bookworm@sha256:02d2537788213fd8b1ecee22ddc42d29e6ce8c3ffc88f7cf12f8a4d92e318715
WORKDIR /app

COPY go.mod .
COPY go.sum .
RUN GOPROXY=https://goproxy.io go mod download

COPY cmd cmd
COPY pkg pkg

CMD ( GOPROXY=https://goproxy.io go test ./... && echo "done" ) || echo "failed"
