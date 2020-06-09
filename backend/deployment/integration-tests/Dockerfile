FROM golang:1.14.1

RUN mkdir /tests
WORKDIR /tests

RUN apt-get update
RUN apt-get install -y bats curl

COPY tests /tests

CMD bats . && echo "done"
