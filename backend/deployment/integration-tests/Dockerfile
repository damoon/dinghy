FROM golang:1.17.8-buster@sha256:0b49d92e9b33cc1972a84942e1a4504038afa73290f73d4b95e00d10e47999f3

RUN mkdir /tests
WORKDIR /tests

RUN apt-get update
RUN apt-get install -y bats curl

COPY deployment/integration-tests/run-tests.sh /usr/local/bin/run-tests.sh
CMD run-tests.sh
