# WORKBENCH ###################################################
FROM golang:1.14.1

RUN mkdir /app
WORKDIR /app

COPY go.mod /app
COPY go.sum /app
RUN go mod download

COPY cmd /app/cmd
COPY pkg /app/pkg

CMD go test ./... && echo "done"
