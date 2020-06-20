#!/bin/bash

set -euxo pipefail;

echo -n "waiting for backend"
while ! curl --max-time 1 --fail --silent -v -o /dev/null http://backend:8090/healthz ;
    do echo -n ".";
done
echo ""

( bats . && echo "done" ) || echo "failed"
