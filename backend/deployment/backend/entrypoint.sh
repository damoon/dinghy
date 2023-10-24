#!/bin/sh
find /go/bin/backend | entr -n -d -r backend server $@
