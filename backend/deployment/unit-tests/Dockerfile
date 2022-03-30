FROM golang:1.18.0-buster@sha256:e3233ebd68cb8b86e22db7765d6376b4f0435e98833493d2ae68e352c12144fe
WORKDIR /app

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY cmd cmd
COPY pkg pkg

CMD ( go test ./... && echo "done" ) || echo "failed"
