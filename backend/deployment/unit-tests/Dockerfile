FROM golang:1.17.6-buster@sha256:2f0dfcbc013e1b3ff28b61d89ed605f7d302815a6ea68675992dd4630d62b916
WORKDIR /app

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY cmd cmd
COPY pkg pkg

CMD ( go test ./... && echo "done" ) || echo "failed"
