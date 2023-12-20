FROM golang:1.21.5-bookworm@sha256:4891844757a81153b512c91f1af30c6de9edf637e1bbebb571b43a631a861b0c

RUN mkdir /tests
WORKDIR /tests

RUN apt-get update
RUN apt-get install -y bats curl

COPY deployment/integration-tests/run-tests.sh /usr/local/bin/run-tests.sh
CMD run-tests.sh
