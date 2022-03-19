FROM golang:1.18.0-buster@sha256:0938f182601066c14a9a3d371d6c3c20237bbfd487ef3f4479a181a26410a0ef
WORKDIR /app

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY cmd cmd
COPY pkg pkg

CMD ( go test ./... && echo "done" ) || echo "failed"
