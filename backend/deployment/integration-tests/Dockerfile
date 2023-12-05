FROM golang:1.21.5-bookworm@sha256:a6b787c7f9046e3fdaa97bca1f76fd23ff4108f612de885e1af87e0dccc02f99

RUN mkdir /tests
WORKDIR /tests

RUN apt-get update
RUN apt-get install -y bats curl

COPY deployment/integration-tests/run-tests.sh /usr/local/bin/run-tests.sh
CMD run-tests.sh
