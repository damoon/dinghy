FROM golang:1.22.2-bookworm@sha256:8e4358fd54fe1c833af23b9a0b84eb5d41fc3dc67d23fde72cf3de5f6319b445

RUN mkdir /tests
WORKDIR /tests

RUN apt-get update
RUN apt-get install -y bats curl

COPY deployment/integration-tests/run-tests.sh /usr/local/bin/run-tests.sh
CMD run-tests.sh
