FROM golang:1.22.2-bookworm@sha256:3451eec734d0d3d87ea6f240d6854dd71573801157b29fbd2adc359fd0fc714c

RUN mkdir /tests
WORKDIR /tests

RUN apt-get update
RUN apt-get install -y bats curl

COPY deployment/integration-tests/run-tests.sh /usr/local/bin/run-tests.sh
CMD run-tests.sh
