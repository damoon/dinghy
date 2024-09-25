# build-env ###################################################
FROM golang:1.22.3-alpine@sha256:f1fe698725f6ed14eb944dc587591f134632ed47fc0732ec27c7642adbe90618 AS build-env
WORKDIR /app

RUN apk add --update --no-cache ca-certificates entr mailcap

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
FROM alpine:3.20.3@sha256:beefdbd8a1da6d2915566fde36db9db0b524eb737fc57cd1367effd16dc0d06d

ENV JAEGER_SERVICE_NAME Dinghy
ENV JAEGER_SAMPLER_TYPE const
ENV JAEGER_SAMPLER_PARAM 1.0
ENV JAEGER_REPORTER_LOG_SPANS true

RUN apk add --update --no-cache ca-certificates

COPY --from=build-env /etc/mime.types /etc/mime.types
COPY --from=build-env /go/bin/backend /usr/local/bin/backend
ENTRYPOINT [ "backend", "server" ]
