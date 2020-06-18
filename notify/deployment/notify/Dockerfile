# WORKBENCH ###################################################
FROM golang:1.14.1 AS workbench
WORKDIR /app

RUN apt-get update
RUN apt-get install -y entr

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY cmd cmd
COPY pkg pkg
RUN go install ./cmd/notify

ENV JAEGER_SERVICE_NAME Dinghy
ENV JAEGER_SAMPLER_TYPE const
ENV JAEGER_SAMPLER_PARAM 1.0
ENV JAEGER_REPORTER_LOG_SPANS true
ENV JAEGER_AGENT_HOST jaeger
ENV JAEGER_AGENT_PORT 6831

CMD find "/go/bin/notify" | entr -d -r notify


# ARTIFACT ####################################################
FROM alpine AS deploy
COPY --from=workbench /go/bin/notify /notify
ENTRYPOINT [ "/notify" ]
