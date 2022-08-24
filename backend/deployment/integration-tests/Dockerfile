FROM golang:1.19.0-buster@sha256:f72495ebf48dffef2ba69e1df2d2545beb3d3019d9dc2c26bed396b36c38e91d

RUN mkdir /tests
WORKDIR /tests

RUN apt-get update
RUN apt-get install -y bats curl

COPY deployment/integration-tests/run-tests.sh /usr/local/bin/run-tests.sh
CMD run-tests.sh
