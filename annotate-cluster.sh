#!/bin/sh

nodes=$(kubectl get nodes -o go-template --template='{{range .items}}{{printf "%s\n" .metadata.name}}{{end}}')
for node in $nodes; do
  kubectl annotate node "${node}" \
    --overwrite=true \
    tilt.dev/registry=registry.guinan.31j.de \
    tilt.dev/registry-from-cluster=registry.registry.svc:5000
done
