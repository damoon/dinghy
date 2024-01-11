FROM golang:1.21.6-bookworm@sha256:1e8ea75170ea0c4653ceab6d37dd12644d1bb44347d1bae457e5d0475d1d26c4
WORKDIR /app

COPY go.mod .
COPY go.sum .
RUN GOPROXY=https://goproxy.io go mod download

COPY cmd cmd
COPY pkg pkg

CMD ( GOPROXY=https://goproxy.io go test ./... && echo "done" ) || echo "failed"
