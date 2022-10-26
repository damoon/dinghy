FROM golang:1.19.2-buster@sha256:202b7098b6bd5c91faacce750514a61e2dd2c31a93b419eca3bb0af3c409835b

RUN mkdir /tests
WORKDIR /tests

RUN apt-get update
RUN apt-get install -y bats curl

COPY deployment/integration-tests/run-tests.sh /usr/local/bin/run-tests.sh
CMD run-tests.sh
