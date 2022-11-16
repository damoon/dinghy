FROM golang:1.19.3-buster@sha256:93ca306de4e8eca74b80a49fe271efc520001b6eacc0905951cf7f0295e40ecb
WORKDIR /app

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY cmd cmd
COPY pkg pkg

CMD ( go test ./... && echo "done" ) || echo "failed"
