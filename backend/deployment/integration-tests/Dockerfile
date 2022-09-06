FROM golang:1.19.1-buster@sha256:d71125b5d1e299a8d9829eb09dbeb0cf576c461cbded210a9a6242d64072adc9

RUN mkdir /tests
WORKDIR /tests

RUN apt-get update
RUN apt-get install -y bats curl

COPY deployment/integration-tests/run-tests.sh /usr/local/bin/run-tests.sh
CMD run-tests.sh
