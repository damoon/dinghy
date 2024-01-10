FROM golang:1.21.6-bookworm@sha256:688ad7f282ce711cc28676531f45d5b840790e28a689fbd4e8e4ac181e23e9c4
WORKDIR /app

COPY go.mod .
COPY go.sum .
RUN GOPROXY=https://goproxy.io go mod download

COPY cmd cmd
COPY pkg pkg

CMD ( GOPROXY=https://goproxy.io go test ./... && echo "done" ) || echo "failed"
