
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

.PHONY: unit-tests
unit-tests: ##@qa Run unit tests.
	go test ./... -short

.PHONY: tests
tests: ##@qa Run unit and integration tests.
	go test ./... -s3Endpoint=play.min.io -s3AccessID=Q3AM3UQ867SPQQA43P2F -s3AccessKey=zuf+tfteSlswRu7BJ86wekitnifILbZam1KYY3TG -s3Bucket=webdav-to-s3

.PHONY: tests-local
tests-local: ##@qa Run unit and integration tests locally.
	go test ./... -s3Endpoint=127.0.0.1:9000 -s3UseSSL=false -s3AccessID=minio -s3AccessKey=minio123 -s3Bucket=webdav-to-s3
