.PHONY: help
.DEFAULT_GOAL := help

install: ## install project dependency packages
	glide install

run.kafka: ## run kafka for testing
	docker-compose -f kafka/docker-compose.yml up -d

stop.kafka: ## stop test kafka
	docker-compose -f kafka/docker-compose.yml down

test: ## run testcases
	go test -v ./...

help:
	@grep -E '^[a-zA-Z0-9._-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
