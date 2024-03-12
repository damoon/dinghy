FROM golang:1.22.1-bookworm@sha256:7f9c05839800653366fa503d21d48f245f1e60a5eb60043efa14b562fea442b6

RUN mkdir /tests
WORKDIR /tests

RUN apt-get update
RUN apt-get install -y bats curl

COPY deployment/integration-tests/run-tests.sh /usr/local/bin/run-tests.sh
CMD run-tests.sh
