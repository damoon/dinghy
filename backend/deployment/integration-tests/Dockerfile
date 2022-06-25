FROM golang:1.18.3-buster@sha256:3ad5048c0c63cbc32d1c576bda513ebfd6e1914b6cd57837e0cd0a31fe81de69

RUN mkdir /tests
WORKDIR /tests

RUN apt-get update
RUN apt-get install -y bats curl

COPY deployment/integration-tests/run-tests.sh /usr/local/bin/run-tests.sh
CMD run-tests.sh
