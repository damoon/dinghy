FROM golang:1.18.0-buster@sha256:1893d1126ad939d129d86a9dd8824f69548aac3ae9610ea8d45dbfc78bd5c8be

RUN mkdir /tests
WORKDIR /tests

RUN apt-get update
RUN apt-get install -y bats curl

COPY deployment/integration-tests/run-tests.sh /usr/local/bin/run-tests.sh
CMD run-tests.sh
