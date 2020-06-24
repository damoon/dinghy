#!/bin/sh
find /go/bin/backend | entr -d -r backend server $@
