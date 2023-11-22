FROM golang:1.21.4-bookworm@sha256:211553a2281d0d49ffd7d8c766624bc5428e40840a3238cbdd3c868b1d9819b0

RUN mkdir /tests
WORKDIR /tests

RUN apt-get update
RUN apt-get install -y bats curl

COPY deployment/integration-tests/run-tests.sh /usr/local/bin/run-tests.sh
CMD run-tests.sh
