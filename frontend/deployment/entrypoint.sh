#!/usr/bin/env sh

set -e

cd /usr/share/nginx/html/

if test -f "index-production.html"; then
    cat index.html.tpl | envsubst > index.html
    unlink index.html.tpl
fi

exec /docker-entrypoint-original.sh "$@"
