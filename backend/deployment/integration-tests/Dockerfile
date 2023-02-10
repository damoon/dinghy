FROM golang:1.20.0-buster@sha256:216b408bb2196155bc9ce7ed3811e4d16db00ff60fc10fb2b5cd0048052594aa

RUN mkdir /tests
WORKDIR /tests

RUN apt-get update
RUN apt-get install -y bats curl

COPY deployment/integration-tests/run-tests.sh /usr/local/bin/run-tests.sh
CMD run-tests.sh
