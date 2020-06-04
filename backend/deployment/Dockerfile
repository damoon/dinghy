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

RUN mkdir /secret
RUN echo "minio123" > /secret/minio

CMD find /app | entr -r -d /app/backend --s3-endpoint=minio:9000 --s3-access-key=minio --s3-secret-access-key-file=/secret/minio --s3-ssl=false --s3-bucket=dinghy --frontend-url=http://127.0.0.1:8000


# ARTIFACT ####################################################
FROM alpine AS deploy
ENV ENVIRONMENT production
WORKDIR /app

COPY --from=workbench /app/backend /backend

ENTRYPOINT [ "/backend" ]
