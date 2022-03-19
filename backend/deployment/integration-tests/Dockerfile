FROM golang:1.18.0-buster@sha256:a22f80e34a9f21a7e11a8ef8d29e04366da6993bf4f8efa41202899250130e6c

RUN mkdir /tests
WORKDIR /tests

RUN apt-get update
RUN apt-get install -y bats curl

COPY deployment/integration-tests/run-tests.sh /usr/local/bin/run-tests.sh
CMD run-tests.sh
