FROM golang:1.20.2-buster@sha256:97c3e1d5d0b59e92e9e067c1e4192d452566b60512203bf885f568e8e364155a

RUN mkdir /tests
WORKDIR /tests

RUN apt-get update
RUN apt-get install -y bats curl

COPY deployment/integration-tests/run-tests.sh /usr/local/bin/run-tests.sh
CMD run-tests.sh
