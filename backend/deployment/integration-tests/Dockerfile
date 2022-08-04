FROM golang:1.19.0-buster@sha256:a7a23f1fba8390b1e038f017c85259c878f406301643653ec6e5b97e75668789

RUN mkdir /tests
WORKDIR /tests

RUN apt-get update
RUN apt-get install -y bats curl

COPY deployment/integration-tests/run-tests.sh /usr/local/bin/run-tests.sh
CMD run-tests.sh
