FROM golang:1.18.1-buster@sha256:279e4ec837160ba80a6ce81e8ae5a5fa4723dc720e61e9b262b82df2c9b53b0f

RUN mkdir /tests
WORKDIR /tests

RUN apt-get update
RUN apt-get install -y bats curl

COPY deployment/integration-tests/run-tests.sh /usr/local/bin/run-tests.sh
CMD run-tests.sh
