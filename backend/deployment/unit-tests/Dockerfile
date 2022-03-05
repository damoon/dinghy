FROM golang:1.17.8-buster@sha256:267cd0d6449f47031b4a6b386e5652ee572ced4ee7e2ec846352b0132bf0d3e5
WORKDIR /app

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY cmd cmd
COPY pkg pkg

CMD ( go test ./... && echo "done" ) || echo "failed"
