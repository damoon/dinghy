FROM golang:1.18.3-buster@sha256:8ec5e989db0ae2aaeb0015cb8cb0a868d262d188c15bca06f61436c6f2ce5559

RUN mkdir /tests
WORKDIR /tests

RUN apt-get update
RUN apt-get install -y bats curl

COPY deployment/integration-tests/run-tests.sh /usr/local/bin/run-tests.sh
CMD run-tests.sh
