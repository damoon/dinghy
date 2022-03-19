FROM golang:1.18.0-buster@sha256:0938f182601066c14a9a3d371d6c3c20237bbfd487ef3f4479a181a26410a0ef

RUN mkdir /tests
WORKDIR /tests

RUN apt-get update
RUN apt-get install -y bats curl

COPY deployment/integration-tests/run-tests.sh /usr/local/bin/run-tests.sh
CMD run-tests.sh
