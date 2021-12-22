FROM golang:1.17.5-buster@sha256:62a929434e43b91ba7a215c1233d97d0649f0d6c3dc2cb03d0c24ac0f81f1982

RUN mkdir /tests
WORKDIR /tests

RUN apt-get update
RUN apt-get install -y bats curl

COPY deployment/integration-tests/run-tests.sh /usr/local/bin/run-tests.sh
CMD run-tests.sh
