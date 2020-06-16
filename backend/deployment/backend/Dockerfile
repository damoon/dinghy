# WORKBENCH ###################################################
FROM golang:1.14.1 AS workbench
WORKDIR /app

RUN apt-get update
RUN apt-get install -y entr

COPY backend/go.mod /app
COPY backend/go.sum /app
RUN go mod download

COPY backend/cmd /app/cmd
COPY backend/pkg /app/pkg
RUN go install ./cmd/backend

RUN mkdir /secret
RUN echo "minio123" > /secret/minio

ENV JAEGER_SERVICE_NAME Dinghy
ENV JAEGER_SAMPLER_TYPE const
ENV JAEGER_SAMPLER_PARAM 1.0
ENV JAEGER_REPORTER_LOG_SPANS true
ENV JAEGER_AGENT_HOST jaeger
ENV JAEGER_AGENT_PORT 6831

CMD find "/go/bin/backend" | entr -d -r backend --s3-endpoint=minio:9000 --s3-access-key=minio --s3-secret-access-key-file=/secret/minio --s3-ssl=false --s3-bucket=dinghy --frontend-url=http://frontend:8000


# ARTIFACT ####################################################
FROM alpine AS deploy
COPY --from=workbench /go/bin/backend /backend
ENTRYPOINT [ "/backend" ]
