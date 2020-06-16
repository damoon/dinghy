# WORKBENCH ###################################################
FROM golang:1.14.1

RUN mkdir /app
WORKDIR /app

COPY backend/go.mod /app
COPY backend/go.sum /app
RUN go mod download

COPY backend/cmd /app/cmd
COPY backend/pkg /app/pkg

CMD go test ./... && echo "done"
