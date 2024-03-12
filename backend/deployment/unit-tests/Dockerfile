FROM golang:1.22.1-bookworm@sha256:7f9c05839800653366fa503d21d48f245f1e60a5eb60043efa14b562fea442b6
WORKDIR /app

COPY go.mod .
COPY go.sum .
RUN GOPROXY=https://goproxy.io go mod download

COPY cmd cmd
COPY pkg pkg

CMD ( GOPROXY=https://goproxy.io go test ./... && echo "done" ) || echo "failed"
