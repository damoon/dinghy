FROM golang:1.17.8-buster@sha256:267cd0d6449f47031b4a6b386e5652ee572ced4ee7e2ec846352b0132bf0d3e5

RUN mkdir /tests
WORKDIR /tests

RUN apt-get update
RUN apt-get install -y bats curl

COPY deployment/integration-tests/run-tests.sh /usr/local/bin/run-tests.sh
CMD run-tests.sh
