FROM golang:1.18.2-buster@sha256:2f489a35bbc88ca3f64471cc4e6425b57f678ff55a32884a3fb9935408af4f67

RUN mkdir /tests
WORKDIR /tests

RUN apt-get update
RUN apt-get install -y bats curl

COPY deployment/integration-tests/run-tests.sh /usr/local/bin/run-tests.sh
CMD run-tests.sh
