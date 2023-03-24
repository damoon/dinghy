FROM golang:1.20.2-buster@sha256:97c3e1d5d0b59e92e9e067c1e4192d452566b60512203bf885f568e8e364155a
WORKDIR /app

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY cmd cmd
COPY pkg pkg

CMD ( go test ./... && echo "done" ) || echo "failed"
