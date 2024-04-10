FROM golang:1.22.2-bookworm@sha256:3451eec734d0d3d87ea6f240d6854dd71573801157b29fbd2adc359fd0fc714c
WORKDIR /app

COPY go.mod .
COPY go.sum .
RUN GOPROXY=https://goproxy.io go mod download

COPY cmd cmd
COPY pkg pkg

CMD ( GOPROXY=https://goproxy.io go test ./... && echo "done" ) || echo "failed"
