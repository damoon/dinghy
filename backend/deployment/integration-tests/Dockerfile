FROM golang:1.18.1-buster@sha256:6ed1375f000a9086969581dd1dcc676495599228fa6b27c3219c4a76b8f09cc4

RUN mkdir /tests
WORKDIR /tests

RUN apt-get update
RUN apt-get install -y bats curl

COPY deployment/integration-tests/run-tests.sh /usr/local/bin/run-tests.sh
CMD run-tests.sh
