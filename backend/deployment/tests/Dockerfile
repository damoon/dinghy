FROM bats/bats:v1.1.0

RUN apk add --update curl

RUN mkdir /tests
WORKDIR /tests

COPY tests /tests

ENTRYPOINT [ "bats", "." ]
