FROM golang:1.19.0-buster@sha256:e4851ced580e41c72a3933f3bd70165672105fe5dd82031b49c6700de265253b
WORKDIR /app

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY cmd cmd
COPY pkg pkg

CMD ( go test ./... && echo "done" ) || echo "failed"
