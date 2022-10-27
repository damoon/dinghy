FROM golang:1.19.2-buster@sha256:b4480899915b13370161a7e574174212b21baccade35b4df86358d76f83da7f9

RUN mkdir /tests
WORKDIR /tests

RUN apt-get update
RUN apt-get install -y bats curl

COPY deployment/integration-tests/run-tests.sh /usr/local/bin/run-tests.sh
CMD run-tests.sh
