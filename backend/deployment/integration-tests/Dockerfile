FROM golang:1.21.6-bookworm@sha256:d8c365d63879c2312e332cb796961f2695dd65124ceb3c0247d9c5426b7dde5f

RUN mkdir /tests
WORKDIR /tests

RUN apt-get update
RUN apt-get install -y bats curl

COPY deployment/integration-tests/run-tests.sh /usr/local/bin/run-tests.sh
CMD run-tests.sh
