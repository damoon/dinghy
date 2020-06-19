#!/bin/sh
echo "/go/bin/backend" | entr -d -r backend $@
