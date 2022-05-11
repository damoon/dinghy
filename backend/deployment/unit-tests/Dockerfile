FROM golang:1.18.2-buster@sha256:0c60ee5eaeb19e585dd845774cff6cc339b3c73454ed49b350172b811475caac
WORKDIR /app

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY cmd cmd
COPY pkg pkg

CMD ( go test ./... && echo "done" ) || echo "failed"
