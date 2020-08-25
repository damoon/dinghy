# Dinghy

This server forwards PUT and GET requests to S3 (minio).

It creates thumbnails of pictures and extracts archives.

Changes are send in realtime to the browser. 

## Contribute

Set up local host names:

``` bash
echo 127.0.0.1 frontend backend minio | sudo tee --append /etc/hosts > /dev/null
```

Forward remote docker in docker:

``` bash
kubectl apply -f hack/docker-daemon.yaml
kubectl wait --for=condition=available --timeout=600s deployment/docker-in-docker
kubectl port-forward svc/dind 12375:2375
```

Start up deployment:

``` bash
DOCKER_HOST=tcp://127.0.0.1:12375 tilt up
```
