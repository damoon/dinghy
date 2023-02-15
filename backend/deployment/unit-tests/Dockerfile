FROM golang:1.20.1-buster@sha256:4e9578db8d2c7876e85d2435bfb360c10485bafce544cb7206ba27a9a2ca2360
WORKDIR /app

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY cmd cmd
COPY pkg pkg

CMD ( go test ./... && echo "done" ) || echo "failed"
