FROM golang:1.20.3-buster@sha256:6cafa6d5d8e5d33e78470db4851e13739431c069f6ca5357ec235e342350f2ed
WORKDIR /app

COPY go.mod .
COPY go.sum .
RUN GOPROXY=https://goproxy.io go mod download

COPY cmd cmd
COPY pkg pkg

CMD ( GOPROXY=https://goproxy.io go test ./... && echo "done" ) || echo "failed"
