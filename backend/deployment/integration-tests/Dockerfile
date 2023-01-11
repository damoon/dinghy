FROM golang:1.19.5-buster@sha256:33d767d032a989ba380e6448b512e9b7056a29a3efbf0409427fdd8445f392e5

RUN mkdir /tests
WORKDIR /tests

RUN apt-get update
RUN apt-get install -y bats curl

COPY deployment/integration-tests/run-tests.sh /usr/local/bin/run-tests.sh
CMD run-tests.sh
