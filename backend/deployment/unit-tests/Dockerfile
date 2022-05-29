FROM golang:1.18.2-buster@sha256:2f489a35bbc88ca3f64471cc4e6425b57f678ff55a32884a3fb9935408af4f67
WORKDIR /app

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY cmd cmd
COPY pkg pkg

CMD ( go test ./... && echo "done" ) || echo "failed"
