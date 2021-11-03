FROM golang:1.17.2-buster@sha256:9cf49bf1e8cf0c8d995899546bea5d876e59b5a9fdeb1575cb6197f97973c20e

RUN mkdir /tests
WORKDIR /tests

RUN apt-get update
RUN apt-get install -y bats curl

COPY deployment/integration-tests/run-tests.sh /usr/local/bin/run-tests.sh
CMD run-tests.sh
