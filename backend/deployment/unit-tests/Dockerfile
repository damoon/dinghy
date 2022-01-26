FROM golang:1.17.6-buster@sha256:4f847cb80c04274125bb3a5af71bf236e6480395667f4621e9f9007cef62464b
WORKDIR /app

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY cmd cmd
COPY pkg pkg

CMD ( go test ./... && echo "done" ) || echo "failed"
