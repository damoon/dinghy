FROM golang:1.19.1-buster@sha256:7656e693167ae561f883fe292ad3e27e05eb681bbfb53ebe0fbb31c95cbc4967

RUN mkdir /tests
WORKDIR /tests

RUN apt-get update
RUN apt-get install -y bats curl

COPY deployment/integration-tests/run-tests.sh /usr/local/bin/run-tests.sh
CMD run-tests.sh
