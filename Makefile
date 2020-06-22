
cli-frontend:
	docker run -ti --rm -v $(shell realpath ./frontend):/app -w /app $(shell docker build --quiet --target build-env -f frontend/deployment/Dockerfile frontend) bash

cli-backend:
	docker run -ti --rm -v $(shell realpath ./backend):/app -w /app $(shell docker build --quiet --target build-env -f backend/deployment/backend/Dockerfile backend) bash

lint: ##@qa run linting for golang.
	docker run -ti --rm -v $(shell realpath ./backend):/app -w /app golangci/golangci-lint:v1.27.0-alpine golangci-lint run --enable-all ./...
