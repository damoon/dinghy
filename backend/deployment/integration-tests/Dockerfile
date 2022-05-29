FROM golang:1.18.2-buster@sha256:13465f0f4ad5bc05ce6ea281c921e5ab577600859f863a6bfaec2ec4d519556f

RUN mkdir /tests
WORKDIR /tests

RUN apt-get update
RUN apt-get install -y bats curl

COPY deployment/integration-tests/run-tests.sh /usr/local/bin/run-tests.sh
CMD run-tests.sh
