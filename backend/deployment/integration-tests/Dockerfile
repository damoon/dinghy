FROM golang:1.17.6-buster@sha256:d3ca64c1e652d7beee33930b0aabb2ab668097690b66698208551e700dbb2e3a

RUN mkdir /tests
WORKDIR /tests

RUN apt-get update
RUN apt-get install -y bats curl

COPY deployment/integration-tests/run-tests.sh /usr/local/bin/run-tests.sh
CMD run-tests.sh
