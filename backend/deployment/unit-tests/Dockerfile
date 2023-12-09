FROM golang:1.21.5-bookworm@sha256:2d3b13c2a6368032e9697e64e7c923184e6e3be03cf01eadff27de124114e64e
WORKDIR /app

COPY go.mod .
COPY go.sum .
RUN GOPROXY=https://goproxy.io go mod download

COPY cmd cmd
COPY pkg pkg

CMD ( GOPROXY=https://goproxy.io go test ./... && echo "done" ) || echo "failed"
