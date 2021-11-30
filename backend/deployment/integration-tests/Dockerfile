FROM golang:1.17.3-buster@sha256:622a414b498370198e16b23f2d174b8aefcd1fd96958b7b6e875dc37bea114ba

RUN mkdir /tests
WORKDIR /tests

RUN apt-get update
RUN apt-get install -y bats curl

COPY deployment/integration-tests/run-tests.sh /usr/local/bin/run-tests.sh
CMD run-tests.sh
