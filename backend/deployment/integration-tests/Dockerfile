FROM golang:1.19.3-buster@sha256:0092b239c5b5dc572729d8fd3c1842513ed1222502f86be65cbb875c5b9b1e8a

RUN mkdir /tests
WORKDIR /tests

RUN apt-get update
RUN apt-get install -y bats curl

COPY deployment/integration-tests/run-tests.sh /usr/local/bin/run-tests.sh
CMD run-tests.sh
