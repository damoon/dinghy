FROM golang:1.19.0-buster@sha256:f72495ebf48dffef2ba69e1df2d2545beb3d3019d9dc2c26bed396b36c38e91d
WORKDIR /app

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY cmd cmd
COPY pkg pkg

CMD ( go test ./... && echo "done" ) || echo "failed"
