FROM golang:1.17.5-buster@sha256:007c4b51431057912f538a180a486e06c972e5462dce1fffbbd3e1ba229949d6

RUN mkdir /tests
WORKDIR /tests

RUN apt-get update
RUN apt-get install -y bats curl

COPY deployment/integration-tests/run-tests.sh /usr/local/bin/run-tests.sh
CMD run-tests.sh
