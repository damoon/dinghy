FROM golang:1.17.6-buster@sha256:a5ba2c081bc7174317ce7f54cf5689c2365454b11a5e271a911d6a2edec1163d

RUN mkdir /tests
WORKDIR /tests

RUN apt-get update
RUN apt-get install -y bats curl

COPY deployment/integration-tests/run-tests.sh /usr/local/bin/run-tests.sh
CMD run-tests.sh
