FROM golang:1.18.2-buster@sha256:a66be0e47d87a1f0626aa8ab7f550d607091e52266fb14a3776615f6fa82d186

RUN mkdir /tests
WORKDIR /tests

RUN apt-get update
RUN apt-get install -y bats curl

COPY deployment/integration-tests/run-tests.sh /usr/local/bin/run-tests.sh
CMD run-tests.sh
