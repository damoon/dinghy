FROM golang:1.21.6-bookworm@sha256:688ad7f282ce711cc28676531f45d5b840790e28a689fbd4e8e4ac181e23e9c4

RUN mkdir /tests
WORKDIR /tests

RUN apt-get update
RUN apt-get install -y bats curl

COPY deployment/integration-tests/run-tests.sh /usr/local/bin/run-tests.sh
CMD run-tests.sh
