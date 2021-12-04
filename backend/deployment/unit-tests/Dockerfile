FROM golang:1.17.4-buster@sha256:e0bac40798d319e8c1cdd8e829bc95061412f5fc5e29fe5e8c857379e5aec52a
WORKDIR /app

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY cmd cmd
COPY pkg pkg

CMD ( go test ./... && echo "done" ) || echo "failed"
