FROM golang:1.17.5-buster@sha256:35b01a3e4adf4c4a403d0ef6879ed3f6e628d17e0be0b543cec89bd480ca4482

RUN mkdir /tests
WORKDIR /tests

RUN apt-get update
RUN apt-get install -y bats curl

COPY deployment/integration-tests/run-tests.sh /usr/local/bin/run-tests.sh
CMD run-tests.sh
