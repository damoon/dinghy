FROM golang:1.18.2-buster@sha256:13465f0f4ad5bc05ce6ea281c921e5ab577600859f863a6bfaec2ec4d519556f
WORKDIR /app

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY cmd cmd
COPY pkg pkg

CMD ( go test ./... && echo "done" ) || echo "failed"
