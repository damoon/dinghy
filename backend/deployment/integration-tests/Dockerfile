FROM golang:1.21.5-bookworm@sha256:f67ff020a28389f545ab57baa063c2650c7dbc3eb3ded0611d7633c713d8433b

RUN mkdir /tests
WORKDIR /tests

RUN apt-get update
RUN apt-get install -y bats curl

COPY deployment/integration-tests/run-tests.sh /usr/local/bin/run-tests.sh
CMD run-tests.sh
