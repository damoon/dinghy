
cli-frontend:
	docker run -ti --rm -v $(shell realpath ./frontend):/app -w /app $(shell docker build --quiet frontend/deployment) bash

cli-backend:
	docker run -ti --rm -v $(shell realpath ./backend):/app -w /app $(shell docker build --quiet backend/deployment) bash

lint: ##@qa run linting for golang.
	docker run -ti --rm -v $(shell realpath ./backend):/app -w /app golangci/golangci-lint:v1.27.0-alpine golangci-lint run --enable-all ./...
