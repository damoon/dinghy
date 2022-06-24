FROM golang:1.18.3-buster@sha256:8c085a2e3269d0d6eff26ca3a628d54808059f5a945bcae08a7e3a8dbcff132b
WORKDIR /app

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY cmd cmd
COPY pkg pkg

CMD ( go test ./... && echo "done" ) || echo "failed"
