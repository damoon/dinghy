FROM golang:1.17.3-buster@sha256:61dfba01ef3c1e3192fb0ee4d4638a7e9143e7056f2235268d269582f803f725
WORKDIR /app

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY cmd cmd
COPY pkg pkg

CMD ( go test ./... && echo "done" ) || echo "failed"
