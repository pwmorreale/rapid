CMD := rapid
TARGET := target

.PHONY: all build clean test lint release goimports tests debug

all: build

help: ## Display this help screen
	@grep -h -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'


tests: lint test race staticcheck ## Run all tests

generate:  ## Generate API and test mock code
	@printf "\033[36m%-30s\033[0m %s\n" "### make $@"
	@go generate ./...

lint:  ## Lint the files
	@printf "\033[36m%-30s\033[0m %s\n" "### make $@"
	@revive -set_exit_status ./...

# use bash to return proper return value from colorize pipe
test: SHELL = /bin/bash
test: .SHELLFLAGS = -o pipefail -c

test:  ## Run unittests
	@printf "\033[36m%-30s\033[0m %s\n" "### make $@"
	@go test -v -vet=all -cover ./... | sed ''/PASS/s//$$(printf "\033[32mPASS\033[0m")/'' | sed ''/FAIL/s//$$(printf "\033[31mFAIL\033[0m")/''


race:  ## Run data race detector
	@printf "\033[36m%-30s\033[0m %s\n" "### make $@"
	@go test -race -short ./...

staticcheck: ## Run staticcheck
	@printf "\033[36m%-30s\033[0m %s\n" "### make $@"
	@staticcheck -f stylish ./...

build:  ## Build the midp microservice.
	@printf "\033[36m%-30s\033[0m %s\n" "### make $@"
	go build -o $(TARGET)/$(CMD)  ./main.go

release: ## Docker Build
	@printf "\033[36m%-30s\033[0m %s\n" "### make $@"
	podman build -t $(IMAGE_REPOSITORY) \
	  $(DOCKER_BUILD_ARGS) \
	  .
	docker push $(IMAGE_REPOSITORY)

clean: ## Remove previous build
	rm -rf $(TARGET)

coverage: ## Display test coverage
	@go test -vet=all -coverprofile=/tmp/coverage.out ./...
	@go tool cover -html=/tmp/coverage.out
