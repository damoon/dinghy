FROM golang:1.18.1-buster@sha256:d8f864bac466e488bb425a59ff9bc104828cd44d0985a8d2fcb5570a9250258a

RUN mkdir /tests
WORKDIR /tests

RUN apt-get update
RUN apt-get install -y bats curl

COPY deployment/integration-tests/run-tests.sh /usr/local/bin/run-tests.sh
CMD run-tests.sh
