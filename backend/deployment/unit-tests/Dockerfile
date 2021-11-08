FROM golang:1.17.3-buster@sha256:eac5f1bc63ff9a91539f0048afc710e690623f75995e79a5da515004ea8de826
WORKDIR /app

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY cmd cmd
COPY pkg pkg

CMD ( go test ./... && echo "done" ) || echo "failed"
