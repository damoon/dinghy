FROM golang:1.18.0-buster@sha256:35cf8e5001bd5b2c3fc0286e379bed24cc0dd1b5b5d50ddd524b30e09696b0cc

RUN mkdir /tests
WORKDIR /tests

RUN apt-get update
RUN apt-get install -y bats curl

COPY deployment/integration-tests/run-tests.sh /usr/local/bin/run-tests.sh
CMD run-tests.sh
