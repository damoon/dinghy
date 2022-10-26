FROM golang:1.19.2-buster@sha256:d7067103aa14d258d61c6ac1c4903779653fb43211346f267f2bf2e86d0212fc

RUN mkdir /tests
WORKDIR /tests

RUN apt-get update
RUN apt-get install -y bats curl

COPY deployment/integration-tests/run-tests.sh /usr/local/bin/run-tests.sh
CMD run-tests.sh
