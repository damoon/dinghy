FROM golang:1.17.6-buster@sha256:babb3d0646633747024a5e39d23ccd01da07957009956a0c5d11144e1a4a8360
WORKDIR /app

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY cmd cmd
COPY pkg pkg

CMD ( go test ./... && echo "done" ) || echo "failed"
