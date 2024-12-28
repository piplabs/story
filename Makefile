help:  ## Display this help message
	@egrep -h '\s##\s' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m  %-30s\033[0m %s\n", $$1, $$2}'

###############################################################################
###                                Build                                    ###
###############################################################################

.PHONY: build
build: mod ## Build the story client.
	@mkdir -p build/
	@go build -o build/story ./client

.PHONY: mod 
mod: ## Update all go.mod files.
	@go mod tidy

###############################################################################
###                                Contracts                                 ###
###############################################################################

.PHONY: contracts-bindings
contract-bindings: ## Generate golang contract bindings.
	make -C ./contracts bindings

###############################################################################
###                                Utils                                 	###
###############################################################################

.PHONY: ensure-detect-secrets
ensure-detect-secrets: ## Checks if detect-secrets is installed.
	@which detect-secrets > /dev/null || echo "detect-secrets not installed, see https://github.com/Yelp/detect-secrets?tab=readme-ov-file#installation"

.PHONY: install-pre-commit
install-pre-commit: ## Installs the pre-commit tool as the git pre-commit hook for this repo.
	@which pre-commit > /dev/null || echo "pre-commit not installed, see https://pre-commit.com/#install"
	@pre-commit install --install-hooks

.PHONY: install-go-tools
install-go-tools: ## Installs the go-dev-tools, like buf.
	@go generate scripts/tools.go

.PHONY: lint
lint: ## Runs linters via pre-commit.
	@pre-commit run -v --all-files

.PHONY: bufgen
bufgen: ## Generates protobufs using buf generate.
	@./scripts/protocgen.sh

.PHONY:
secrets-baseline: ensure-detect-secrets ## Update secrets baseline.
	@detect-secrets scan --exclude-file pnpm-lock.yaml > .secrets.baseline

.PHONY: fix-golden
fix-golden: ## Fixes golden test fixtures.
	@./scripts/fix_golden_tests.sh

###############################################################################
###                                Testing                                 	###
###############################################################################
