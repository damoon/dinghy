FROM golang:1.18.5-buster@sha256:cb94fa4e6f97e8ce76a1809e4e1d071a39753c1bc696035a1992b9dc3d444ebc

RUN mkdir /tests
WORKDIR /tests

RUN apt-get update
RUN apt-get install -y bats curl

COPY deployment/integration-tests/run-tests.sh /usr/local/bin/run-tests.sh
CMD run-tests.sh
