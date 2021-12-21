# build-env ###################################################
FROM golang:1.17.5-buster@sha256:35b01a3e4adf4c4a403d0ef6879ed3f6e628d17e0be0b543cec89bd480ca4482 AS build-env
WORKDIR /app

RUN apt-get update
RUN apt-get install -y entr

ENV JAEGER_SERVICE_NAME Dinghy
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

COPY cmd cmd
COPY pkg pkg
RUN go install -ldflags="-X main.gitRef=${SOURCE_BRANCH} -X main.gitHash=${SOURCE_COMMIT}" ./cmd/backend

###############################################################
FROM debian:buster@sha256:94ccfd1c5115a6903cbb415f043a0b04e307be3f37b768cf6d6d3edff0021da3

ENV JAEGER_SERVICE_NAME Dinghy
ENV JAEGER_SAMPLER_TYPE const
ENV JAEGER_SAMPLER_PARAM 1.0
ENV JAEGER_REPORTER_LOG_SPANS true

RUN apt-get update && apt-get install -y \
    ca-certificates \
 && rm -rf /var/lib/apt/lists/*

COPY --from=build-env /etc/mime.types /etc/mime.types
COPY --from=build-env /go/bin/backend /usr/local/bin/backend
ENTRYPOINT [ "backend", "server" ]
