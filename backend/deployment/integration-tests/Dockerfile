FROM golang:1.17.7-buster@sha256:270086796f4fea7905b19889d33a6063e93eb6ad9c6fde11aaf3b2f82a757bef

RUN mkdir /tests
WORKDIR /tests

RUN apt-get update
RUN apt-get install -y bats curl

COPY deployment/integration-tests/run-tests.sh /usr/local/bin/run-tests.sh
CMD run-tests.sh
