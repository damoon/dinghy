FROM golang:1.22.2-bookworm@sha256:cafb443b26e25774290119a4fbc393776d8cdba83d999ffc1681f2ef8ccd6869
WORKDIR /app

COPY go.mod .
COPY go.sum .
RUN GOPROXY=https://goproxy.io go mod download

COPY cmd cmd
COPY pkg pkg

CMD ( GOPROXY=https://goproxy.io go test ./... && echo "done" ) || echo "failed"
