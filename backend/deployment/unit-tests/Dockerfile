FROM golang:1.17.5-buster@sha256:007c4b51431057912f538a180a486e06c972e5462dce1fffbbd3e1ba229949d6
WORKDIR /app

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY cmd cmd
COPY pkg pkg

CMD ( go test ./... && echo "done" ) || echo "failed"
