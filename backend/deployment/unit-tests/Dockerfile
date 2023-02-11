FROM golang:1.20.0-buster@sha256:f8fbd744a9129e3b705362d39c17f0cbdafd1979a08508ba2d72e488fdf98e22
WORKDIR /app

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY cmd cmd
COPY pkg pkg

CMD ( go test ./... && echo "done" ) || echo "failed"
