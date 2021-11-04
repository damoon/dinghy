FROM golang:1.17.3-buster@sha256:61dfba01ef3c1e3192fb0ee4d4638a7e9143e7056f2235268d269582f803f725

RUN mkdir /tests
WORKDIR /tests

RUN apt-get update
RUN apt-get install -y bats curl

COPY deployment/integration-tests/run-tests.sh /usr/local/bin/run-tests.sh
CMD run-tests.sh
