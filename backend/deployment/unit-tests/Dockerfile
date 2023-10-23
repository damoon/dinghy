FROM golang:1.20.5-buster@sha256:eb3f9ac805435c1b2c965d63ce460988e1000058e1f67881324746362baf9572
WORKDIR /app

COPY go.mod .
COPY go.sum .
RUN GOPROXY=https://goproxy.io go mod download

COPY cmd cmd
COPY pkg pkg

CMD ( GOPROXY=https://goproxy.io go test ./... && echo "done" ) || echo "failed"
