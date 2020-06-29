#!/usr/bin/env sh

set -e

cd /usr/share/nginx/html/

if test -f "index-production.html"; then
    cat index-production.html | envsubst > index.html
    unlink index-production.html
fi

exec /docker-entrypoint-original.sh "$@"
