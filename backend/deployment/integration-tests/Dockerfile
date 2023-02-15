FROM golang:1.20.1-buster@sha256:4e9578db8d2c7876e85d2435bfb360c10485bafce544cb7206ba27a9a2ca2360

RUN mkdir /tests
WORKDIR /tests

RUN apt-get update
RUN apt-get install -y bats curl

COPY deployment/integration-tests/run-tests.sh /usr/local/bin/run-tests.sh
CMD run-tests.sh
