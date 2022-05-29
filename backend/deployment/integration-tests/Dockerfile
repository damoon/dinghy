FROM golang:1.18.2-buster@sha256:e5013e70182658129a37b368718be0d12deff03655a1514854ac9b1d54d8d94c

RUN mkdir /tests
WORKDIR /tests

RUN apt-get update
RUN apt-get install -y bats curl

COPY deployment/integration-tests/run-tests.sh /usr/local/bin/run-tests.sh
CMD run-tests.sh
