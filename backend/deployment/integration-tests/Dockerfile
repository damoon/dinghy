FROM golang:1.19.4-buster@sha256:c1e438a538be3fe72468f1569a57565225fcab11762f372f46b20d65bba45ab9

RUN mkdir /tests
WORKDIR /tests

RUN apt-get update
RUN apt-get install -y bats curl

COPY deployment/integration-tests/run-tests.sh /usr/local/bin/run-tests.sh
CMD run-tests.sh
