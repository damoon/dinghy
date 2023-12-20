FROM golang:1.21.5-bookworm@sha256:4891844757a81153b512c91f1af30c6de9edf637e1bbebb571b43a631a861b0c
WORKDIR /app

COPY go.mod .
COPY go.sum .
RUN GOPROXY=https://goproxy.io go mod download

COPY cmd cmd
COPY pkg pkg

CMD ( GOPROXY=https://goproxy.io go test ./... && echo "done" ) || echo "failed"
