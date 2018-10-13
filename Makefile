PROJECT_NAME := "gopxe"
PKG := "github.com/ppetko/$(PROJECT_NAME)"
PKG_LIST := $(shell go list ${PKG}/... | grep -v /vendor/)
GO_FILES := $(shell find . -name '*.go' | grep -v /vendor/ | grep -v _test.go)

.PHONY: all dep build docker-build clean test coverage lint

all: build

lint: ## Lint the files
	@golint -set_exit_status ${PKG_LIST}

test: ## Run unittests
	@go test -v ./...

race: dep ## Run data race detector
	@go test -race -short ${PKG_LIST}

cover: ## Generate global code coverage report in HTML
	./tools/coverage.sh html;

dep: ## Get the dependencies
	@go get -v -d ./...

build: dep ## Build the binary file
	@go build -v

docker-build: ## Build docker container from Dockerfile
	@docker build -t $(PROJECT_NAME) -f Dockerfile . 

docker-test: ## Run all go tests in Docker
	@docker run --rm -ti $(PROJECT_NAME) go tests -v ./... 

clean: ## Remove previous build
	@rm -f $(PROJECT_NAME)
	@rm -f coverage.html

help: ## Display this help screen
	@grep -h -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

