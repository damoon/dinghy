FROM golang:1.21.5-bookworm@sha256:1415bb0b25d3bffc0a44dcf9851c20a9f8bbe558095221d931f2e4a4cc3596eb

RUN mkdir /tests
WORKDIR /tests

RUN apt-get update
RUN apt-get install -y bats curl

COPY deployment/integration-tests/run-tests.sh /usr/local/bin/run-tests.sh
CMD run-tests.sh
