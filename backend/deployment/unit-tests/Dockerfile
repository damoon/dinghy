FROM golang:1.22.2-bookworm@sha256:8e4358fd54fe1c833af23b9a0b84eb5d41fc3dc67d23fde72cf3de5f6319b445
WORKDIR /app

COPY go.mod .
COPY go.sum .
RUN GOPROXY=https://goproxy.io go mod download

COPY cmd cmd
COPY pkg pkg

CMD ( GOPROXY=https://goproxy.io go test ./... && echo "done" ) || echo "failed"
