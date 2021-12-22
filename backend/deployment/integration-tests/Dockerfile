FROM golang:1.17.5-buster@sha256:47030a15ca7cccf40da9f968fabebc5e99ce5455ad4a688cda3d748278a9c967

RUN mkdir /tests
WORKDIR /tests

RUN apt-get update
RUN apt-get install -y bats curl

COPY deployment/integration-tests/run-tests.sh /usr/local/bin/run-tests.sh
CMD run-tests.sh
