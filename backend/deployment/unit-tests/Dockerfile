FROM golang:1.19.3-buster@sha256:256a00b2d3966bc34e3664ca09325b2eeb66130ae090cea25658824df04b6163
WORKDIR /app

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY cmd cmd
COPY pkg pkg

CMD ( go test ./... && echo "done" ) || echo "failed"
