FROM golang:1.19.2-buster@sha256:535a9d6a746bfdb02e6c2b506c62aaeb72b8327bcc7e17a5f2cc831a6b4026c2

RUN mkdir /tests
WORKDIR /tests

RUN apt-get update
RUN apt-get install -y bats curl

COPY deployment/integration-tests/run-tests.sh /usr/local/bin/run-tests.sh
CMD run-tests.sh
