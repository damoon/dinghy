FROM golang:1.18.3-buster@sha256:0383eb7a4e40ddbafb2cf875f20e523d10d4a8dc748c474b130dd607de9cb445
WORKDIR /app

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY cmd cmd
COPY pkg pkg

CMD ( go test ./... && echo "done" ) || echo "failed"
