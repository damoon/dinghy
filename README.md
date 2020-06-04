# Webdav to S3 Server

This server forwards PUT and GET requests to S3 (minio).

## Example

start minio

``` bash
make minio
```

start the server

``` bash
make start
```

run the tests

``` bash
make test
```


`kubectl -n container-image-builder port-forward svc/dind 12375:2375`

`DOCKER_HOST=tcp://127.0.0.1:12375 tilt up`

`GO111MODULE=on code .`
