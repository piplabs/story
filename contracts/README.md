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

## Deploy

These smart contracts are predeploys (part of the genesis state of Execution Layer).

To generate this first state:

1. Add a .env file in `contracts/.env`

```
ADMIN_ADDRESS=0x123...
UPGRADE_ADMIN_ADDRESS=0x234...
```
- `ADMIN_ADDRESS` will be the owner of `IPTokenStaking` and `UpgradeEntryPoint`, able to execute admin methods.
- `UPGRADE_ADMIN_ADDRESS` will be the owner of the [ProxyAdmin](https://github.com/OpenZeppelin/openzeppelin-contracts/blob/master/contracts/proxy/transparent/ProxyAdmin.sol) for each upgradeable predeploy.

2. Run
```
forge script script/GenerateAlloc.s.sol -vvvv --chain-id <DESIRED_CHAIN_ID>
```

Copy the contents of the resulting JSON file, and paste in the `alloc` item of `story-geth` `genesis.json`
