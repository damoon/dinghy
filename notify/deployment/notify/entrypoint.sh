#!/bin/sh
echo "/go/bin/notify" | entr -d -r notify $@
