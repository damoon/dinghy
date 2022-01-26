FROM golang:1.17.6-buster@sha256:babb3d0646633747024a5e39d23ccd01da07957009956a0c5d11144e1a4a8360

RUN mkdir /tests
WORKDIR /tests

RUN apt-get update
RUN apt-get install -y bats curl

COPY deployment/integration-tests/run-tests.sh /usr/local/bin/run-tests.sh
CMD run-tests.sh
