FROM golang:1.18.2-buster@sha256:11d168eae2d70b24f0926b73126fe739f91646d6f82e21d1b56dd5d571d8838b

RUN mkdir /tests
WORKDIR /tests

RUN apt-get update
RUN apt-get install -y bats curl

COPY deployment/integration-tests/run-tests.sh /usr/local/bin/run-tests.sh
CMD run-tests.sh
