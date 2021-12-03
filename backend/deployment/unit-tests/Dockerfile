FROM golang:1.17.3-buster@sha256:ee3a388a872237ddb600de3ab9512e73df0043f8878f0f355baeb5b723ef16ec
WORKDIR /app

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY cmd cmd
COPY pkg pkg

CMD ( go test ./... && echo "done" ) || echo "failed"
