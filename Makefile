.PHONY: build
build:
	@go build -o bin/stress-test .

.PHONY: run
run:
	@go run main.go --url=http://google.com --requests=100 --concurrency=10

.PHONY: docker-build
docker-build:
	@docker build -t goexpert-stress-test .

.PHONY: docker-run
docker-run: docker-build
	@docker run --rm goexpert-stress-test --url=http://google.com --requests=100 --concurrency=10

.PHONY: docker-run-custom
docker-run-custom: docker-build
	@docker run --rm goexpert-stress-test --url=$(URL) --requests=$(REQUESTS) --concurrency=$(CONCURRENCY)

.PHONY: test-basic
test-basic:
	@docker compose run --rm test-basic

.PHONY: test-high-load
test-high-load:
	@docker compose run --rm test-high-load

.PHONY: test-low-concurrency
test-low-concurrency:
	@docker compose run --rm test-low-concurrency

.PHONY: test-many-requests
test-many-requests:
	@docker compose run --rm test-many-requests

.PHONY: test-timeout
test-timeout:
	@docker compose run --rm test-timeout

.PHONY: test-errors
test-errors:
	@docker compose run --rm test-errors

.PHONY: test-connection-errors
test-connection-errors:
	@docker compose run --rm test-connection-errors

.PHONY: help-test
help-test:
	@echo "Available test scenarios:"
	@echo "  make test-basic             # 100 requests with 10 concurrent"
	@echo "  make test-high-load         # 1000 requests with 50 concurrent"
	@echo "  make test-low-concurrency   # 100 requests with 1 concurrent"
	@echo "  make test-many-requests     # 10000 requests with 100 concurrent"
	@echo "  make test-timeout           # Test with delayed responses"
	@echo "  make test-errors            # Test with error responses"
	@echo "  make test-connection-errors # Test with connection errors"
