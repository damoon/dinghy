FROM golang:1.17.6-buster@sha256:6ba8e82f63c644170940fa7329b8112c7e56ba3b55812baf364fdc01790196a3

RUN mkdir /tests
WORKDIR /tests

RUN apt-get update
RUN apt-get install -y bats curl

COPY deployment/integration-tests/run-tests.sh /usr/local/bin/run-tests.sh
CMD run-tests.sh
