FROM golang:1.18.3-buster@sha256:8c085a2e3269d0d6eff26ca3a628d54808059f5a945bcae08a7e3a8dbcff132b

RUN mkdir /tests
WORKDIR /tests

RUN apt-get update
RUN apt-get install -y bats curl

COPY deployment/integration-tests/run-tests.sh /usr/local/bin/run-tests.sh
CMD run-tests.sh
