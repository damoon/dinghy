FROM golang:1.21.4-bookworm@sha256:211553a2281d0d49ffd7d8c766624bc5428e40840a3238cbdd3c868b1d9819b0
WORKDIR /app

COPY go.mod .
COPY go.sum .
RUN GOPROXY=https://goproxy.io go mod download

COPY cmd cmd
COPY pkg pkg

CMD ( GOPROXY=https://goproxy.io go test ./... && echo "done" ) || echo "failed"
