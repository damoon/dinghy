FROM golang:1.18.5-buster@sha256:49ae1460441117f4737e4e4bfdc2ff166156abce55e402b501516158afeba36e

RUN mkdir /tests
WORKDIR /tests

RUN apt-get update
RUN apt-get install -y bats curl

COPY deployment/integration-tests/run-tests.sh /usr/local/bin/run-tests.sh
CMD run-tests.sh
