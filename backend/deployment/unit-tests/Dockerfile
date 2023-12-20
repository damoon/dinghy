FROM golang:1.21.5-bookworm@sha256:24e832041c0df9265032b46ce2eb4a5e887119089b84705565f6cf4bd534caf6
WORKDIR /app

COPY go.mod .
COPY go.sum .
RUN GOPROXY=https://goproxy.io go mod download

COPY cmd cmd
COPY pkg pkg

CMD ( GOPROXY=https://goproxy.io go test ./... && echo "done" ) || echo "failed"
