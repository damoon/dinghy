FROM golang:1.17.3-buster@sha256:6d0772539cd83e5a8b24a0003e82e042148f89762d36c610c220536322541e94

RUN mkdir /tests
WORKDIR /tests

RUN apt-get update
RUN apt-get install -y bats curl

COPY deployment/integration-tests/run-tests.sh /usr/local/bin/run-tests.sh
CMD run-tests.sh
