FROM golang:1.18.2-buster@sha256:30968aa42e8f6f3cf45d12b22ade8aad14d68183ceec5c310e07c01ce49d7f98

RUN mkdir /tests
WORKDIR /tests

RUN apt-get update
RUN apt-get install -y bats curl

COPY deployment/integration-tests/run-tests.sh /usr/local/bin/run-tests.sh
CMD run-tests.sh
