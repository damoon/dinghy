FROM golang:1.17.7-buster@sha256:3d11eb4c0032307bbecbb40c30c0719006f4960f8ae4d01383e6b442425a9057

RUN mkdir /tests
WORKDIR /tests

RUN apt-get update
RUN apt-get install -y bats curl

COPY deployment/integration-tests/run-tests.sh /usr/local/bin/run-tests.sh
CMD run-tests.sh
