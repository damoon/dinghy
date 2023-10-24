#!/bin/sh
find /go/bin/notify | entr -n -d -r notify server $@
