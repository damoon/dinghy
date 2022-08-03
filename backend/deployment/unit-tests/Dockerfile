FROM golang:1.18.5-buster@sha256:49ae1460441117f4737e4e4bfdc2ff166156abce55e402b501516158afeba36e
WORKDIR /app

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY cmd cmd
COPY pkg pkg

CMD ( go test ./... && echo "done" ) || echo "failed"
