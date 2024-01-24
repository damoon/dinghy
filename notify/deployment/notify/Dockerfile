# build-env ###################################################
FROM golang:1.21.6-alpine@sha256:51a7800206bc7b276a9d62a7229cdede7b1e0f45ec28259ed44c1603c6cda1e7 AS build-env
WORKDIR /app

RUN apk add --update --no-cache ca-certificates entr

ENV JAEGER_SERVICE_NAME Dinghy-Notify
ENV JAEGER_SAMPLER_TYPE const
ENV JAEGER_SAMPLER_PARAM 1.0
ENV JAEGER_REPORTER_LOG_SPANS true
ENV JAEGER_AGENT_HOST jaeger
ENV JAEGER_AGENT_PORT 6831

COPY deployment/notify/entrypoint.sh /entrypoint.sh
ENTRYPOINT ["/entrypoint.sh"]

RUN GRPC_HEALTH_PROBE_VERSION=v0.4.21 && \
    wget -qO/bin/grpc_health_probe https://github.com/grpc-ecosystem/grpc-health-probe/releases/download/${GRPC_HEALTH_PROBE_VERSION}/grpc_health_probe-linux-amd64 && \
    chmod +x /bin/grpc_health_probe

COPY go.mod .
COPY go.sum .
RUN GOPROXY=https://goproxy.io go mod download

ARG SOURCE_BRANCH
ARG SOURCE_COMMIT

COPY cmd cmd
COPY pkg pkg
RUN go install -ldflags="-X main.gitRef=${SOURCE_BRANCH} -X main.gitHash=${SOURCE_COMMIT}" ./cmd/notify

###############################################################
FROM alpine:3.19.0@sha256:51b67269f354137895d43f3b3d810bfacd3945438e94dc5ac55fdac340352f48

ENV JAEGER_SERVICE_NAME Dinghy-Notify
ENV JAEGER_SAMPLER_TYPE const
ENV JAEGER_SAMPLER_PARAM 1.0
ENV JAEGER_REPORTER_LOG_SPANS true

RUN apk add --update --no-cache ca-certificates

COPY --from=build-env /bin/grpc_health_probe /bin/grpc_health_probe
COPY --from=build-env /go/bin/notify /usr/local/bin/notify
ENTRYPOINT [ "notify", "server" ]
