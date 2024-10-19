ifneq ("$(wildcard .env)","")
	include .env
	export $(shell sed 's/=.*//' .env)
endif


.PHONY: help
help:  ## Display this help message.
	@egrep -h '\s##\s' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m  %-30s\033[0m %s\n", $$1, $$2}'

.PHONY: version
version: ## Print tool versions.
	@forge --version
	@abigen --version

.PHONY: build
build: version ## Build contracts.
	forge build

.PHONY: test
test: version ## Run tests.
	forge test

CONTRACTS := IPTokenStaking UpgradeEntrypoint UBIPool Create3

.PHONY: bindings
bindings: check-abigen-version build ## Generate golang contract bindings.
	./bindings/scripts/gen.sh $(CONTRACTS)
	./bindings/scripts/genmore.sh $(CONTRACTS)

.PHONY: fork-holesky
fork-holesky: ## Run an anvil holesky fork.
	anvil --fork-url https://holesky.infura.io/v3/$(INFURA_KEY)

.PHONY: fork-mainnet
fork-mainnet: ## Run an anvil mainnet fork.
	anvil --fork-url https://mainnet.infura.io/v3/$(INFURA_KEY)

.PHONY: check-abigen-version
check-abigen-version: ## Check abigen version, exit(1) if not 1.13.14-stable.
	@version=$$(abigen --version); \
	if [ "$$version" != "abigen version 1.13.14-stable" ]; then \
		echo "abigen version is not 1.13.14-stable"; \
		echo "Install with go install github.com/ethereum/go-ethereum/cmd/abigen@v1.13.14"; \
		exit 1; \
	fi

.PHONY: gengen
gengen: ## Generate genesis config
	forge script script/GenerateAlloc.s.sol -vvv --chain-id $(CHAIN_ID)
