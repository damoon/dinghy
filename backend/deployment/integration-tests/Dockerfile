FROM golang:1.18.4-buster@sha256:477b10a289d7be9c5124da52f2ff86d13c24e591b82da3b6b1ad3b8f1a04c8d9

RUN mkdir /tests
WORKDIR /tests

RUN apt-get update
RUN apt-get install -y bats curl

COPY deployment/integration-tests/run-tests.sh /usr/local/bin/run-tests.sh
CMD run-tests.sh
