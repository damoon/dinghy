FROM golang:1.22.2-bookworm@sha256:cafb443b26e25774290119a4fbc393776d8cdba83d999ffc1681f2ef8ccd6869

RUN mkdir /tests
WORKDIR /tests

RUN apt-get update
RUN apt-get install -y bats curl

COPY deployment/integration-tests/run-tests.sh /usr/local/bin/run-tests.sh
CMD run-tests.sh
