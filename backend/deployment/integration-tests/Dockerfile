FROM golang:1.19.3-buster@sha256:93ca306de4e8eca74b80a49fe271efc520001b6eacc0905951cf7f0295e40ecb

RUN mkdir /tests
WORKDIR /tests

RUN apt-get update
RUN apt-get install -y bats curl

COPY deployment/integration-tests/run-tests.sh /usr/local/bin/run-tests.sh
CMD run-tests.sh
