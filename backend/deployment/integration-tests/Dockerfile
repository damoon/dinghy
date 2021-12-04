FROM golang:1.17.4-buster@sha256:e0bac40798d319e8c1cdd8e829bc95061412f5fc5e29fe5e8c857379e5aec52a

RUN mkdir /tests
WORKDIR /tests

RUN apt-get update
RUN apt-get install -y bats curl

COPY deployment/integration-tests/run-tests.sh /usr/local/bin/run-tests.sh
CMD run-tests.sh
