FROM golang:1.17.3-buster@sha256:eac5f1bc63ff9a91539f0048afc710e690623f75995e79a5da515004ea8de826

RUN mkdir /tests
WORKDIR /tests

RUN apt-get update
RUN apt-get install -y bats curl

COPY deployment/integration-tests/run-tests.sh /usr/local/bin/run-tests.sh
CMD run-tests.sh
