FROM golang:1.19.2-buster@sha256:60477d9f1c5829330227946c1e58f402f7a06d4626eed0dc621687bd3b2d4cf5

RUN mkdir /tests
WORKDIR /tests

RUN apt-get update
RUN apt-get install -y bats curl

COPY deployment/integration-tests/run-tests.sh /usr/local/bin/run-tests.sh
CMD run-tests.sh
