FROM golang:1.18.4-buster@sha256:295d17b21f1b206054dcd72fa1a1b26f197fc661b0901a312ee471e40b251757

RUN mkdir /tests
WORKDIR /tests

RUN apt-get update
RUN apt-get install -y bats curl

COPY deployment/integration-tests/run-tests.sh /usr/local/bin/run-tests.sh
CMD run-tests.sh
