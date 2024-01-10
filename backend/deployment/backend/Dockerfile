# build-env ###################################################
FROM golang:1.21.6-alpine@sha256:fd78f2fb1e49bcf343079bbbb851c936a18fc694df993cbddaa24ace0cc724c5 AS build-env
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
FROM alpine:3.19.0@sha256:51b67269f354137895d43f3b3d810bfacd3945438e94dc5ac55fdac340352f48

ENV JAEGER_SERVICE_NAME Dinghy
ENV JAEGER_SAMPLER_TYPE const
ENV JAEGER_SAMPLER_PARAM 1.0
ENV JAEGER_REPORTER_LOG_SPANS true

RUN apk add --update --no-cache ca-certificates

COPY --from=build-env /etc/mime.types /etc/mime.types
COPY --from=build-env /go/bin/backend /usr/local/bin/backend
ENTRYPOINT [ "backend", "server" ]
