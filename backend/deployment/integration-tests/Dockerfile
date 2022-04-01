FROM golang:1.18.0-buster@sha256:284da7e6ad7fd1250a785c24912fb6f8652176b588d637de5c5b3d4a4aa1693a

RUN mkdir /tests
WORKDIR /tests

RUN apt-get update
RUN apt-get install -y bats curl

COPY deployment/integration-tests/run-tests.sh /usr/local/bin/run-tests.sh
CMD run-tests.sh
