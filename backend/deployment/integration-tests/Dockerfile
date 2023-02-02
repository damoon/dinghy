FROM golang:1.20.0-buster@sha256:fc1e0b02674cef154ceca0ba6ca4d3658bd385d5faadecbbd44809727a56329d

RUN mkdir /tests
WORKDIR /tests

RUN apt-get update
RUN apt-get install -y bats curl

COPY deployment/integration-tests/run-tests.sh /usr/local/bin/run-tests.sh
CMD run-tests.sh
