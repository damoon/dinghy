# WORKBENCH ###################################################
FROM golang:1.14.1 AS workbench
ENV ENVIRONMENT development
WORKDIR /app

RUN apt-get update
RUN apt-get install -y entr

COPY go.mod /app
COPY go.sum /app
RUN go mod download

COPY cmd /app/cmd
COPY pkg /app/pkg
RUN go build ./cmd/backend

CMD find /app/backend /app/web | entr -r -d /app/backend


# ARTIFACT ####################################################
FROM alpine AS deploy
ENV ENVIRONMENT production
WORKDIR /app

COPY --from=workbench /app/backend /app/backend

CMD /app/backend
