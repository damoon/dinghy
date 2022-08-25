FROM golang:1.19.0-buster@sha256:d84495e2ad12dfeb51dec3c9f170659ebd09db234d6990b3364f877785684a14

RUN mkdir /tests
WORKDIR /tests

RUN apt-get update
RUN apt-get install -y bats curl

COPY deployment/integration-tests/run-tests.sh /usr/local/bin/run-tests.sh
CMD run-tests.sh
