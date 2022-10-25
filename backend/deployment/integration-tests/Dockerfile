FROM golang:1.19.2-buster@sha256:afa7136268fc9a2719b9e170fdc5df2d3d0b120fbefabbf7f1b24e124a3b3808

RUN mkdir /tests
WORKDIR /tests

RUN apt-get update
RUN apt-get install -y bats curl

COPY deployment/integration-tests/run-tests.sh /usr/local/bin/run-tests.sh
CMD run-tests.sh
