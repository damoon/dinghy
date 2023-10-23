# build-env ###################################################
FROM golang:1.20.3-buster@sha256:73c225bc5e2353f20dbe0466819b70a51a114a93bfe4af035a3bb9e1ecdd4107 AS build-env
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
RUN GOPROXY=https://goproxy.io go mod download

ARG SOURCE_BRANCH
ARG SOURCE_COMMIT

COPY cmd cmd
COPY pkg pkg
RUN go install -ldflags="-X main.gitRef=${SOURCE_BRANCH} -X main.gitHash=${SOURCE_COMMIT}" ./cmd/backend

###############################################################
FROM debian:bullseye@sha256:c141beaa9e0767774221cc82efe3a6712a1cc4f75d2699334dfd9a28a6f7357b

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
