FROM golang:1.18.2-buster@sha256:11d168eae2d70b24f0926b73126fe739f91646d6f82e21d1b56dd5d571d8838b
WORKDIR /app

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY cmd cmd
COPY pkg pkg

CMD ( go test ./... && echo "done" ) || echo "failed"
