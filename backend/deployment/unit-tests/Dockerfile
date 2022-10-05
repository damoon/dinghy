FROM golang:1.19.2-buster@sha256:60477d9f1c5829330227946c1e58f402f7a06d4626eed0dc621687bd3b2d4cf5
WORKDIR /app

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY cmd cmd
COPY pkg pkg

CMD ( go test ./... && echo "done" ) || echo "failed"
