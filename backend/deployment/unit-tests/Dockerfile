FROM golang:1.22.3-bookworm@sha256:5c56bd47228dd572d8a82971cf1f946cd8bb1862a8ec6dc9f3d387cc94136976
WORKDIR /app

COPY go.mod .
COPY go.sum .
RUN GOPROXY=https://goproxy.io go mod download

COPY cmd cmd
COPY pkg pkg

CMD ( GOPROXY=https://goproxy.io go test ./... && echo "done" ) || echo "failed"
