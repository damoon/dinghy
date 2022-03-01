FROM golang:1.17.7-buster@sha256:428a77349bc4b83a13ad950160288e323a6552680c11a7473f776cbd97beae6b

RUN mkdir /tests
WORKDIR /tests

RUN apt-get update
RUN apt-get install -y bats curl

COPY deployment/integration-tests/run-tests.sh /usr/local/bin/run-tests.sh
CMD run-tests.sh
