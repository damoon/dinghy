FROM golang:1.21.5-bookworm@sha256:f67ff020a28389f545ab57baa063c2650c7dbc3eb3ded0611d7633c713d8433b
WORKDIR /app

COPY go.mod .
COPY go.sum .
RUN GOPROXY=https://goproxy.io go mod download

COPY cmd cmd
COPY pkg pkg

CMD ( GOPROXY=https://goproxy.io go test ./... && echo "done" ) || echo "failed"
