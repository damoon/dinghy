FROM golang:1.17.3-buster@sha256:970e472a45917b6bb7a20b9d269216cebadc7fd3afd4ce6e4690636b562afa7c

RUN mkdir /tests
WORKDIR /tests

RUN apt-get update
RUN apt-get install -y bats curl

COPY deployment/integration-tests/run-tests.sh /usr/local/bin/run-tests.sh
CMD run-tests.sh
