FROM golang:1.20.2-buster@sha256:5f9d227a93a4453df44827dbf44acef84e0b0c3abf39470d8a69b28bbd56733f

RUN mkdir /tests
WORKDIR /tests

RUN apt-get update
RUN apt-get install -y bats curl

COPY deployment/integration-tests/run-tests.sh /usr/local/bin/run-tests.sh
CMD run-tests.sh
