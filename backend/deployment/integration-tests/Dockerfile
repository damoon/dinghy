FROM golang:1.21.6-bookworm@sha256:1e8ea75170ea0c4653ceab6d37dd12644d1bb44347d1bae457e5d0475d1d26c4

RUN mkdir /tests
WORKDIR /tests

RUN apt-get update
RUN apt-get install -y bats curl

COPY deployment/integration-tests/run-tests.sh /usr/local/bin/run-tests.sh
CMD run-tests.sh
