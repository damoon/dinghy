FROM golang:1.21.4-bookworm@sha256:690cc6dc987070de71f7dbeaf4d16cbc6f190d99283660a81afbf7baeb06d028
WORKDIR /app

COPY go.mod .
COPY go.sum .
RUN GOPROXY=https://goproxy.io go mod download

COPY cmd cmd
COPY pkg pkg

CMD ( GOPROXY=https://goproxy.io go test ./... && echo "done" ) || echo "failed"
