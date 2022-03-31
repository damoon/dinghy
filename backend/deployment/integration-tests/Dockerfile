FROM golang:1.18.0-buster@sha256:f2b1b8a903a8204cb0e84db2348552df552af477390c2307f75a493138b58da8

RUN mkdir /tests
WORKDIR /tests

RUN apt-get update
RUN apt-get install -y bats curl

COPY deployment/integration-tests/run-tests.sh /usr/local/bin/run-tests.sh
CMD run-tests.sh
