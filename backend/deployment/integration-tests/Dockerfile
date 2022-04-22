FROM golang:1.18.1-buster@sha256:a58a53ecc3088e9ab801e0a8851dae8294f9931446685c43dda9646144f38af3

RUN mkdir /tests
WORKDIR /tests

RUN apt-get update
RUN apt-get install -y bats curl

COPY deployment/integration-tests/run-tests.sh /usr/local/bin/run-tests.sh
CMD run-tests.sh
