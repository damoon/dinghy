FROM golang:1.17.3-buster@sha256:099de80fda2a712d3ee9f66f7f8b62061aad4b5e0ff43f66c86f16ba6a3a89cd

RUN mkdir /tests
WORKDIR /tests

RUN apt-get update
RUN apt-get install -y bats curl

COPY deployment/integration-tests/run-tests.sh /usr/local/bin/run-tests.sh
CMD run-tests.sh
