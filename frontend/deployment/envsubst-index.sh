#!/usr/bin/env sh

set -e

cd /usr/share/nginx/html/

if test -f "index.html.tpl"; then
    mv index.html.tpl ../index.html.tpl
fi

cat ../index.html.tpl | envsubst > index.html
