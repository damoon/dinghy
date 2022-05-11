FROM golang:1.18.2-buster@sha256:0c60ee5eaeb19e585dd845774cff6cc339b3c73454ed49b350172b811475caac

RUN mkdir /tests
WORKDIR /tests

RUN apt-get update
RUN apt-get install -y bats curl

COPY deployment/integration-tests/run-tests.sh /usr/local/bin/run-tests.sh
CMD run-tests.sh
