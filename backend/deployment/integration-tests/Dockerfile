FROM golang:1.19.3-buster@sha256:5d6475d21c6cd8f3a8f578f51b0a87bc198b245898242116a6a2fdae747de47b

RUN mkdir /tests
WORKDIR /tests

RUN apt-get update
RUN apt-get install -y bats curl

COPY deployment/integration-tests/run-tests.sh /usr/local/bin/run-tests.sh
CMD run-tests.sh
