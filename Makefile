
cli-frontend:
	docker run -ti --rm -v $(shell realpath ./frontend):/app -w /app $(shell docker build --quiet frontend/deployment) bash

cli-backend:
	docker run -ti --rm -v $(shell realpath ./backend):/app -w /app $(shell docker build --quiet backend/deployment) bash

lint: ##@qa run linting for golang.
	docker run -ti --rm -v $(shell realpath ./backend):/app -w /app golangci/golangci-lint:v1.27.0-alpine golangci-lint run --enable-all ./...

.PHONY: unit-tests
unit-tests: ##@qa Run unit tests.
	go test ./... -short

.PHONY: tests
tests-local: ##@qa Run unit and integration tests locally.
	go test ./... -s3Endpoint=127.0.0.1:9000 -s3UseSSL=false -s3AccessID=minio -s3AccessKey=minio123 -s3Bucket=webdav-to-s3

generate:
	go generate ./backend/pkg
	go generate ./notify/pkg
