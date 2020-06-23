# build-env ###################################################
FROM golang:1.14.4 AS build-env
WORKDIR /app

RUN apt-get update
RUN apt-get install -y entr

ENV JAEGER_SERVICE_NAME Backend
ENV JAEGER_SAMPLER_TYPE const
ENV JAEGER_SAMPLER_PARAM 1.0
ENV JAEGER_REPORTER_LOG_SPANS true
ENV JAEGER_AGENT_HOST jaeger
ENV JAEGER_AGENT_PORT 6831

COPY deployment/backend/entrypoint.sh /entrypoint.sh
ENTRYPOINT ["/entrypoint.sh"]

COPY go.mod .
COPY go.sum .
RUN go mod download

ARG SOURCE_BRANCH
ARG SOURCE_COMMIT
ENV LDFLAGS="-X main.gitRef=${SOURCE_BRANCH} -X main.gitHash=${SOURCE_COMMIT} "

COPY cmd cmd
COPY pkg pkg
RUN go install -ldflags="${LDFLAGS}" ./cmd/backend

###############################################################
FROM debian:buster

ENV JAEGER_SERVICE_NAME Backend
ENV JAEGER_SAMPLER_TYPE const
ENV JAEGER_SAMPLER_PARAM 1.0
ENV JAEGER_REPORTER_LOG_SPANS true

COPY --from=build-env /etc/mime.types /etc/mime.types
COPY --from=build-env /go/bin/backend /backend
ENTRYPOINT [ "/backend" ]
