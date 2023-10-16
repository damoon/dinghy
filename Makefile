
cli-frontend:
	docker run -ti --rm -v $(shell realpath ./frontend):/app -w /app $(shell docker build --quiet --target build-env -f frontend/deployment/Dockerfile frontend) bash

cli-backend:
	docker run -ti --rm -v $(shell realpath ./backend):/app -w /app $(shell docker build --quiet --target build-env -f backend/deployment/backend/Dockerfile backend) bash

elm-analyse:
	docker run -ti --rm -v $(shell realpath ./frontend):/app -w /app $(shell docker build --quiet --target build-env -f frontend/deployment/Dockerfile frontend) elm-analyse

golangci-lint-backend:
	docker run -ti --rm -v $(shell realpath ./backend):/app -w /app golangci/golangci-lint:v1.43.0-alpine golangci-lint run --modules-download-mode=readonly ./...

golangci-lint-notify:
	docker run -ti --rm -v $(shell realpath ./notify):/app -w /app golangci/golangci-lint:v1.43.0-alpine golangci-lint run --modules-download-mode=readonly ./...
