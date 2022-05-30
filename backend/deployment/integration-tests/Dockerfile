FROM golang:1.18.2-buster@sha256:ad5d8648c50bd3b2dfdb3bd4f6ab52a82b43d949f53c1998884b1653e0b873b5

RUN mkdir /tests
WORKDIR /tests

RUN apt-get update
RUN apt-get install -y bats curl

COPY deployment/integration-tests/run-tests.sh /usr/local/bin/run-tests.sh
CMD run-tests.sh
