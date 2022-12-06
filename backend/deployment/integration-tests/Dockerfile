FROM golang:1.19.3-buster@sha256:256a00b2d3966bc34e3664ca09325b2eeb66130ae090cea25658824df04b6163

RUN mkdir /tests
WORKDIR /tests

RUN apt-get update
RUN apt-get install -y bats curl

COPY deployment/integration-tests/run-tests.sh /usr/local/bin/run-tests.sh
CMD run-tests.sh
