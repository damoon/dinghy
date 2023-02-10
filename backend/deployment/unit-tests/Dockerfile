FROM golang:1.20.0-buster@sha256:216b408bb2196155bc9ce7ed3811e4d16db00ff60fc10fb2b5cd0048052594aa
WORKDIR /app

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY cmd cmd
COPY pkg pkg

CMD ( go test ./... && echo "done" ) || echo "failed"
