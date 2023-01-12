FROM golang:1.19.5-buster@sha256:150a98859da4c88935958d21d4fca1c74a139eac88e1eb2364c8f60d69236ede

RUN mkdir /tests
WORKDIR /tests

RUN apt-get update
RUN apt-get install -y bats curl

COPY deployment/integration-tests/run-tests.sh /usr/local/bin/run-tests.sh
CMD run-tests.sh
