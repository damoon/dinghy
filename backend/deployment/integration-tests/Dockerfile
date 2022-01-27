FROM golang:1.17.6-buster@sha256:2f0dfcbc013e1b3ff28b61d89ed605f7d302815a6ea68675992dd4630d62b916

RUN mkdir /tests
WORKDIR /tests

RUN apt-get update
RUN apt-get install -y bats curl

COPY deployment/integration-tests/run-tests.sh /usr/local/bin/run-tests.sh
CMD run-tests.sh
