FROM golang:1.18.0-buster@sha256:f2b1b8a903a8204cb0e84db2348552df552af477390c2307f75a493138b58da8
WORKDIR /app

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY cmd cmd
COPY pkg pkg

CMD ( go test ./... && echo "done" ) || echo "failed"
