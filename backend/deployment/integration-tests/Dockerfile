FROM golang:1.17.3-buster@sha256:48f6285a02b7682373e97f153238f3a831268ca2e7dd6204bf540b6801ac63bc

RUN mkdir /tests
WORKDIR /tests

RUN apt-get update
RUN apt-get install -y bats curl

COPY deployment/integration-tests/run-tests.sh /usr/local/bin/run-tests.sh
CMD run-tests.sh
