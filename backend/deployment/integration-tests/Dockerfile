FROM golang:1.17.5-buster@sha256:ee9771f71e7f70eb8c09852ab05844d9e896a00f931af40fbb5b9481c732230c

RUN mkdir /tests
WORKDIR /tests

RUN apt-get update
RUN apt-get install -y bats curl

COPY deployment/integration-tests/run-tests.sh /usr/local/bin/run-tests.sh
CMD run-tests.sh
