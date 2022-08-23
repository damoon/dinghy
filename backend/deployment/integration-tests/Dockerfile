FROM golang:1.19.0-buster@sha256:e4851ced580e41c72a3933f3bd70165672105fe5dd82031b49c6700de265253b

RUN mkdir /tests
WORKDIR /tests

RUN apt-get update
RUN apt-get install -y bats curl

COPY deployment/integration-tests/run-tests.sh /usr/local/bin/run-tests.sh
CMD run-tests.sh
