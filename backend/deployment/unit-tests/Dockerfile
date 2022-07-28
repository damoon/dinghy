FROM golang:1.18.4-buster@sha256:135b3520e4ecac1e3d4b2718ad10e8f1e2a84a8c2f5aa4e74e37810bb0b0c301
WORKDIR /app

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY cmd cmd
COPY pkg pkg

CMD ( go test ./... && echo "done" ) || echo "failed"
