# WORKBENCH ###################################################
FROM golang:1.14.1 AS workbench
WORKDIR /app

RUN apt-get update
RUN apt-get install -y entr

COPY notification/go.mod /app
COPY notification/go.sum /app
RUN go mod download

COPY notification/cmd /app/cmd
COPY notification/pkg /app/pkg
RUN go install ./cmd/notification

ENV JAEGER_SERVICE_NAME Dinghy
ENV JAEGER_SAMPLER_TYPE const
ENV JAEGER_SAMPLER_PARAM 1.0
ENV JAEGER_REPORTER_LOG_SPANS true
ENV JAEGER_AGENT_HOST jaeger
ENV JAEGER_AGENT_PORT 6831

CMD find "/go/bin/notification" | entr -d -r notification


# ARTIFACT ####################################################
FROM alpine AS deploy
COPY --from=workbench /go/bin/notification /notification
ENTRYPOINT [ "/notification" ]
