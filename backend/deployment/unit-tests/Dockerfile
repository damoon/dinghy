FROM golang:1.19.1-buster@sha256:4b2498de3f543907f4954584ce12558fc3cf80844907f3e7b9f7e98c3223b949
WORKDIR /app

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY cmd cmd
COPY pkg pkg

CMD ( go test ./... && echo "done" ) || echo "failed"
