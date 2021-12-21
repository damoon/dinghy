FROM golang:1.17.5-buster@sha256:35b01a3e4adf4c4a403d0ef6879ed3f6e628d17e0be0b543cec89bd480ca4482
WORKDIR /app

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY cmd cmd
COPY pkg pkg

CMD ( go test ./... && echo "done" ) || echo "failed"
