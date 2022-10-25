FROM golang:1.19.2-buster@sha256:afa7136268fc9a2719b9e170fdc5df2d3d0b120fbefabbf7f1b24e124a3b3808
WORKDIR /app

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY cmd cmd
COPY pkg pkg

CMD ( go test ./... && echo "done" ) || echo "failed"
