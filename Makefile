
include ./etc/help.mk

lint: ##@qa run linting for golang.
	golangci-lint run --enable-all ./...

.PHONY: minio
minio: ##@development Start minio server (port:9000, user:minio, secret:minio123).
	docker run --rm \
	-p 9000:9000 \
	-e MINIO_ACCESS_KEY=minio \
	-e MINIO_SECRET_KEY=minio123 \
	minio/minio \
	server /data
#	-e MINIO_HTTP_TRACE=/dev/stdout \

start: ##@development Start the server (port:8080, admin port:8081).
	DEBUG=0 air -c air.conf

.PHONY: test
test: ##@qa Run checks.
	go test ./...
	bats tests
