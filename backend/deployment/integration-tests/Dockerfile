FROM golang:1.19.4-buster@sha256:7c97baef6d7a48c2520691e38cca1c56a86b0cad36831c01d959913711f30aa3

RUN mkdir /tests
WORKDIR /tests

RUN apt-get update
RUN apt-get install -y bats curl

COPY deployment/integration-tests/run-tests.sh /usr/local/bin/run-tests.sh
CMD run-tests.sh
