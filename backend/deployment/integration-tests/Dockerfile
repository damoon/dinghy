FROM golang:1.18.2-buster@sha256:8c7ceffb849c06888d9b8d6ee1314883aa9d9f26ab3aa544e7e0dcac3ec41389

RUN mkdir /tests
WORKDIR /tests

RUN apt-get update
RUN apt-get install -y bats curl

COPY deployment/integration-tests/run-tests.sh /usr/local/bin/run-tests.sh
CMD run-tests.sh
