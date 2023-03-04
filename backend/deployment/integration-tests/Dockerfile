FROM golang:1.20.1-buster@sha256:9628a1ad396ad02e6d26b8506cd36b2e04c56f827f060bd06fd095963263fea2

RUN mkdir /tests
WORKDIR /tests

RUN apt-get update
RUN apt-get install -y bats curl

COPY deployment/integration-tests/run-tests.sh /usr/local/bin/run-tests.sh
CMD run-tests.sh
