# Story Contracts

## Install Dependencies
1. Install `npm` if you haven't.

2. Pull `node_modules`.

```
npm install -g pnpm
pnpm install
```

## Build

1. Install `abigen`.

```
go install github.com/ethereum/go-ethereum/cmd/abigen@latest
```

2. Build the contracts.

```
make build
```

## Test

1. Install `foundry`.

```
curl -L https://foundry.paradigm.xyz | bash
source ~/.bash_profile
foundryup
```

2. Run tests.

```
make test
```
