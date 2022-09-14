FROM golang:1.19.1-buster@sha256:4b2498de3f543907f4954584ce12558fc3cf80844907f3e7b9f7e98c3223b949

RUN mkdir /tests
WORKDIR /tests

RUN apt-get update
RUN apt-get install -y bats curl

COPY deployment/integration-tests/run-tests.sh /usr/local/bin/run-tests.sh
CMD run-tests.sh
