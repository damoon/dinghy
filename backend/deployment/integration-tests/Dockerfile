FROM golang:1.21.6-bookworm@sha256:c4b696f1b2bf7d42e02b62b160c3f81c39386e1c567603df8c514ad6ce93361d

RUN mkdir /tests
WORKDIR /tests

RUN apt-get update
RUN apt-get install -y bats curl

COPY deployment/integration-tests/run-tests.sh /usr/local/bin/run-tests.sh
CMD run-tests.sh
