FROM golang:1.18.2-buster@sha256:a66be0e47d87a1f0626aa8ab7f550d607091e52266fb14a3776615f6fa82d186
WORKDIR /app

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY cmd cmd
COPY pkg pkg

CMD ( go test ./... && echo "done" ) || echo "failed"
