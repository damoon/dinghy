FROM golang:1.21.4-bookworm@sha256:690cc6dc987070de71f7dbeaf4d16cbc6f190d99283660a81afbf7baeb06d028

RUN mkdir /tests
WORKDIR /tests

RUN apt-get update
RUN apt-get install -y bats curl

COPY deployment/integration-tests/run-tests.sh /usr/local/bin/run-tests.sh
CMD run-tests.sh
