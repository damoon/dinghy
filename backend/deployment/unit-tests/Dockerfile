FROM golang:1.17.7-buster@sha256:428a77349bc4b83a13ad950160288e323a6552680c11a7473f776cbd97beae6b
WORKDIR /app

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY cmd cmd
COPY pkg pkg

CMD ( go test ./... && echo "done" ) || echo "failed"
