FROM golang:1.21.5-bookworm@sha256:24e832041c0df9265032b46ce2eb4a5e887119089b84705565f6cf4bd534caf6

RUN mkdir /tests
WORKDIR /tests

RUN apt-get update
RUN apt-get install -y bats curl

COPY deployment/integration-tests/run-tests.sh /usr/local/bin/run-tests.sh
CMD run-tests.sh
