FROM golang:1.21.3-bookworm@sha256:deebfdaab16f8508e2590048d416b451e76ebe22e487cafd68990ff28f93ab1c
WORKDIR /app

COPY go.mod .
COPY go.sum .
RUN GOPROXY=https://goproxy.io go mod download

COPY cmd cmd
COPY pkg pkg

CMD ( GOPROXY=https://goproxy.io go test ./... && echo "done" ) || echo "failed"
