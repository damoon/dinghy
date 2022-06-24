FROM golang:1.18.3-buster@sha256:0383eb7a4e40ddbafb2cf875f20e523d10d4a8dc748c474b130dd607de9cb445

RUN mkdir /tests
WORKDIR /tests

RUN apt-get update
RUN apt-get install -y bats curl

COPY deployment/integration-tests/run-tests.sh /usr/local/bin/run-tests.sh
CMD run-tests.sh
