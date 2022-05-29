# build-env ###################################################
FROM golang:1.18.2-buster@sha256:e5013e70182658129a37b368718be0d12deff03655a1514854ac9b1d54d8d94c AS build-env
WORKDIR /app

RUN apt-get update
RUN apt-get install -y entr

ENV JAEGER_SERVICE_NAME Dinghy-Notify
ENV JAEGER_SAMPLER_TYPE const
ENV JAEGER_SAMPLER_PARAM 1.0
ENV JAEGER_REPORTER_LOG_SPANS true
ENV JAEGER_AGENT_HOST jaeger
ENV JAEGER_AGENT_PORT 6831

COPY deployment/notify/entrypoint.sh /entrypoint.sh
ENTRYPOINT ["/entrypoint.sh"]

RUN GRPC_HEALTH_PROBE_VERSION=v0.3.1 && \
    wget -qO/bin/grpc_health_probe https://github.com/grpc-ecosystem/grpc-health-probe/releases/download/${GRPC_HEALTH_PROBE_VERSION}/grpc_health_probe-linux-amd64 && \
    chmod +x /bin/grpc_health_probe

COPY go.mod .
COPY go.sum .
RUN go mod download

ARG SOURCE_BRANCH
ARG SOURCE_COMMIT

COPY cmd cmd
COPY pkg pkg
RUN go install -ldflags="-X main.gitRef=${SOURCE_BRANCH} -X main.gitHash=${SOURCE_COMMIT}" ./cmd/notify

###############################################################
FROM debian:bullseye@sha256:3f1d6c17773a45c97bd8f158d665c9709d7b29ed7917ac934086ad96f92e4510

ENV JAEGER_SERVICE_NAME Dinghy-Notify
ENV JAEGER_SAMPLER_TYPE const
ENV JAEGER_SAMPLER_PARAM 1.0
ENV JAEGER_REPORTER_LOG_SPANS true

COPY --from=build-env /bin/grpc_health_probe /bin/grpc_health_probe
COPY --from=build-env /go/bin/notify /usr/local/bin/notify
ENTRYPOINT [ "notify", "server" ]
