.DEFAULT_GOAL := help

GO := go
LDFLAGS := '-w -s'
BUILD_DIR := ./build
APP_BIN := $(BUILD_DIR)/floridaman
APP_VERSION := 1.0

SOURCE_FILES := $(shell find . -type f -name "*.go")

COVERAGE_FILE := $(BUILD_DIR)/coverage.out

.PHONY: help
help: ## prints help
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

$(APP_BIN): $(SOURCE_FILES)
	@mkdir -p $(BUILD_DIR)
	@CGO_ENABLED=0 $(GO) build -ldflags $(LDFLAGS) -o $(BUILD_DIR)/floridaman ./cmd/floridaman

build: $(APP_BIN) ## builds floridaman's binary for production use

.PHONY: test
test: ## runs all floridaman tests
	@$(GO) test ./...

$(COVERAGE_FILE): $(SOURCE_FILES)
	@$(GO) test -coverprofile=$(COVERAGE_FILE) ./...

coverage: $(COVERAGE_FILE) ## tests floridaman codebase and generates a coverage file

.PHONY: see-coverage
see-coverage: coverage ## showcases the coverage in a browser using html output
	@$(GO) tool cover -html=$(COVERAGE_FILE)

.PHONY: clean
clean: ## removes build assets generated by the build target from the system
	@rm -rf $(BUILD_DIR)/*

.PHONY: lint
lint: ## detects flaws in the code and checks for style
	@docker run --rm -v $$PWD:/app -w /app golangci/golangci-lint:v1.47.2 golangci-lint run

.PHONY: docker
docker: ## builds the application's docker image
	@docker build -t davidonium/floridaman:$(APP_VERSION) .
	@docker tag davidonium/floridaman:$(APP_VERSION) davidonium/floridaman:latest

deploy: build ## deploys the application to production, assumes that the config in ~/.ssh/config has a Host floridaman entry
	# copy the binary generated by build target
	scp $(APP_BIN) floridaman:/opt/floridaman

	# reload the systemd service
	ssh floridaman -c "sudo systemctl reload floridaman"
