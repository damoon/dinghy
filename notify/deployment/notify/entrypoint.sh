#!/bin/sh
find /go/bin/notify | entr -d -r notify server $@
