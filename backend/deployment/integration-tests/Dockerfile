FROM golang:1.18.0-buster@sha256:e3233ebd68cb8b86e22db7765d6376b4f0435e98833493d2ae68e352c12144fe

RUN mkdir /tests
WORKDIR /tests

RUN apt-get update
RUN apt-get install -y bats curl

COPY deployment/integration-tests/run-tests.sh /usr/local/bin/run-tests.sh
CMD run-tests.sh
