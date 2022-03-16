FROM golang:1.18.0-buster@sha256:e4950365f46f107ba92f8a9557defccb38f151a79fd5453e9424ec34de08c35d

RUN mkdir /tests
WORKDIR /tests

RUN apt-get update
RUN apt-get install -y bats curl

COPY deployment/integration-tests/run-tests.sh /usr/local/bin/run-tests.sh
CMD run-tests.sh
