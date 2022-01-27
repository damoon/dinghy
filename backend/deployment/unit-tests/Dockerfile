FROM golang:1.17.6-buster@sha256:6ba8e82f63c644170940fa7329b8112c7e56ba3b55812baf364fdc01790196a3
WORKDIR /app

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY cmd cmd
COPY pkg pkg

CMD ( go test ./... && echo "done" ) || echo "failed"
