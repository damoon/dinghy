FROM golang:1.18.2-buster@sha256:949a195c5d44e7b7f71499c146def008c7e98b8d02c276e567ab2e5c6971587b

RUN mkdir /tests
WORKDIR /tests

RUN apt-get update
RUN apt-get install -y bats curl

COPY deployment/integration-tests/run-tests.sh /usr/local/bin/run-tests.sh
CMD run-tests.sh
