FROM golang:1.17.5-buster@sha256:7f89be6ec02f7eb4b429e9a43f4ac7d49d2c0cb79a99a28e2d1efe1e7e8a3473

RUN mkdir /tests
WORKDIR /tests

RUN apt-get update
RUN apt-get install -y bats curl

COPY deployment/integration-tests/run-tests.sh /usr/local/bin/run-tests.sh
CMD run-tests.sh
