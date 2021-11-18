FROM golang:1.17.3-buster@sha256:29f41d09f4e91ac5c516c0b92610dfc372f8713db81feb24d45b67eb9232d260

RUN mkdir /tests
WORKDIR /tests

RUN apt-get update
RUN apt-get install -y bats curl

COPY deployment/integration-tests/run-tests.sh /usr/local/bin/run-tests.sh
CMD run-tests.sh
