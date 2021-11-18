FROM golang:1.17.3-buster@sha256:970e472a45917b6bb7a20b9d269216cebadc7fd3afd4ce6e4690636b562afa7c
WORKDIR /app

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY cmd cmd
COPY pkg pkg

CMD ( go test ./... && echo "done" ) || echo "failed"
