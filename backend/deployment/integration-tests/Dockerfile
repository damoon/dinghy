FROM golang:1.18.0-buster@sha256:4d2553850f501c598d05e9a587653a5a81d6c7967f4368ade8e492fcf250f8f1

RUN mkdir /tests
WORKDIR /tests

RUN apt-get update
RUN apt-get install -y bats curl

COPY deployment/integration-tests/run-tests.sh /usr/local/bin/run-tests.sh
CMD run-tests.sh
