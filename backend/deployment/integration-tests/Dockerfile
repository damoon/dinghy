FROM golang:1.18.3-buster@sha256:7acdbec4276ce3c8915ffa03148bd4b41dfb46a190fd1242e9df5b11061d282d

RUN mkdir /tests
WORKDIR /tests

RUN apt-get update
RUN apt-get install -y bats curl

COPY deployment/integration-tests/run-tests.sh /usr/local/bin/run-tests.sh
CMD run-tests.sh
