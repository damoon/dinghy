FROM golang:1.18.5-buster@sha256:eb4b3bd59da4604b0ab27d3837e39dacfba4a48cdd289419a3c49993c0eadba4

RUN mkdir /tests
WORKDIR /tests

RUN apt-get update
RUN apt-get install -y bats curl

COPY deployment/integration-tests/run-tests.sh /usr/local/bin/run-tests.sh
CMD run-tests.sh
