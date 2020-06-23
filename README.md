# Dinghy

This server forwards PUT and GET requests to S3 (minio).

It creates thumbnails of pictures and extracts archives.

Changes are send in realtime to the browser. 

## Contribute

Set up local host names:

`echo 127.0.0.1 frontend backend minio | sudo tee --append /etc/hosts > /dev/null`

Forward remote docker in docker:

`kubectl -n container-image-builder port-forward svc/dind 12375:2375`

Start up deployment:

`DOCKER_HOST=tcp://127.0.0.1:12375 tilt up`
