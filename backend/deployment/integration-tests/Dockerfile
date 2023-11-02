FROM golang:1.21.3-bookworm@sha256:02d2537788213fd8b1ecee22ddc42d29e6ce8c3ffc88f7cf12f8a4d92e318715

RUN mkdir /tests
WORKDIR /tests

RUN apt-get update
RUN apt-get install -y bats curl

COPY deployment/integration-tests/run-tests.sh /usr/local/bin/run-tests.sh
CMD run-tests.sh
