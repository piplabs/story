## Installation

Building the node requires a working Go (version 1.22 or higher, see `go.mod`) and `goreleaser` ([see installation guide here](https://goreleaser.com/install/) or install with `make ensure-go-releaser`). You can install them using your favorite package manager. Once the dependencies are installed, run:

```bash
make build-docker
```

## Usage

### Testing

To run the end-to-end tests, run:

```bash
MANIFEST=simple make e2e-run
```

### Starting a devnet

To start a devnet, run:

```bash
make devnet-deploy
```

To stop it, run:

```bash
make devnet-clean
```

### Setup

To setup the environment for development, run:
```bash
go mod download
make ensure-go-releaser
make install-go-tools
make build-docker
make build-iliad
make contracts-bindings
```

## Directory Structure

<pre>
├── <a href="./contracts/">contracts</a>: Solidity contracts, bindings, and testing for the Iliad protocol.
│ ├── <a href="./contracts/bindings/">bindings</a>: Go bindings for smart contracts.
│ ├── <a href="./contracts/src/">src</a>: Solidity source files for the protocol's smart contracts.
│ └── <a href="./contracts/test/">test</a>: Tests for smart contracts.
├── <a href="./docs/">docs</a>: Documentation resources, including images and diagrams.
├── <a href="./client/">iliad</a>: The Iliad instance, including application logic mechanisms.
│ ├── <a href="./client/app/">app</a>: Application logic for Iliad.
│ └── <a href="./client/cmd/">cmd</a>: Command-line tools and utilities.
├── <a href="./lib/">lib</a>: Core libraries for various protocol functionalities.
├── <a href="./scripts/">scripts</a>: Utility scripts for development and operational tasks.
└── <a href="./test/">test</a>: Testing suite for end-to-end, smoke, and utility testing.
</pre>

## Contributing

For detailed instructions on how to contribute, including our coding standards, testing practices, and how to submit pull requests, please see [the contribution guidelines](./docs/contributing.md).

## Acknowledgements

The development of Iliad has been a journey of learning, adaptation, and innovation. As we built Iliad, we drew inspiration and knowledge from the work of several remarkable teams within the blockchain and Ethereum ecosystem.

We extend our gratitude to the following teams for their pioneering work and the open-source resources they've provided, which have significantly influenced our development process:

- [**CometBFT**](https://github.com/cometbft/cometbft): Our heartfelt thanks go to the CometBFT team for their groundbreaking work in consensus algorithms.

- [**Geth**](https://github.com/ethereum/go-ethereum): The go-ethereum team's relentless dedication to the Ethereum ecosystem has been nothing short of inspiring. Their comprehensive and robust implementation of the Ethereum protocol has provided us with a solid foundation to build upon.

- [**Erigon**](https://github.com/ledgerwatch/erigon): We are grateful to Erigon for their novel work in maximizing EVM performance.

- [**Optimism**](https://github.com/ethereum-optimism/optimism): We thank the Optimism team for their dedication to open source work within the Ethereum ecosystem.

Acknowledging these teams is not only a gesture of gratitude but also a reflection of our commitment to collaborative progress and the open-source ethos. The path they've paved has enabled us to contribute our innovations back to the community, and we look forward to continuing this tradition of mutual growth and support.

## Security



Please refer to [SECURITY.md](./SECURITY.md).
