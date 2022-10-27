FROM golang:1.19.2-buster@sha256:b4480899915b13370161a7e574174212b21baccade35b4df86358d76f83da7f9
WORKDIR /app

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY cmd cmd
COPY pkg pkg

CMD ( go test ./... && echo "done" ) || echo "failed"
