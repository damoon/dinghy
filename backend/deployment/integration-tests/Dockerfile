FROM golang:1.17.6-buster@sha256:4f847cb80c04274125bb3a5af71bf236e6480395667f4621e9f9007cef62464b

RUN mkdir /tests
WORKDIR /tests

RUN apt-get update
RUN apt-get install -y bats curl

COPY deployment/integration-tests/run-tests.sh /usr/local/bin/run-tests.sh
CMD run-tests.sh
