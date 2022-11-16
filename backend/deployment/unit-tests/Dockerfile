FROM golang:1.19.3-buster@sha256:5d6475d21c6cd8f3a8f578f51b0a87bc198b245898242116a6a2fdae747de47b
WORKDIR /app

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY cmd cmd
COPY pkg pkg

CMD ( go test ./... && echo "done" ) || echo "failed"
