FROM golang:1.21.3-bookworm@sha256:deebfdaab16f8508e2590048d416b451e76ebe22e487cafd68990ff28f93ab1c

RUN mkdir /tests
WORKDIR /tests

RUN apt-get update
RUN apt-get install -y bats curl

COPY deployment/integration-tests/run-tests.sh /usr/local/bin/run-tests.sh
CMD run-tests.sh
