---
sidebar_position: 1
---

# `x/mint`

## Contents

- [`x/mint`](#xmint)
  - [Contents](#contents)
  - [State](#state)
    - [Params](#params)
  - [Begin-Block](#begin-block)
    - [Inflation amount calculation](#inflation-amount-calculation)
  - [Parameters](#parameters)
  - [Events](#events)
    - [BeginBlocker](#beginblocker)

## State

### Params

The mint module stores its params in state with the prefix of `0x01`,
it can be updated by authority via specific smart contract.

* Params: `mint/params -> legacy_amino(params)`

```protobuf
message Params {
  option (amino.name) = "client/x/mint/Params";

  // type of coin to mint
  string mint_denom = 1;
  // inflation amount per year
  uint64 inflations_per_year = 2;
  // expected blocks per year
  uint64 blocks_per_year = ;
}
```

## Begin-Block

Minting parameters are calculated and inflation paid at the beginning of each block.

### Inflation amount calculation

Inflation amount is calculated using an "inflation calculation function" that's
passed to the `NewAppModule` function. If no function is passed, then the SDK's
default inflation function will be used (`DefaultInflationCalculationFn`). In case a custom
inflation calculation logic is needed, this can be achieved by defining and
passing a function that matches `InflationCalculationFn`'s signature.

```go
type InflationCalculationFn func(ctx sdk.Context, minter Minter, params Params, bondedRatio math.LegacyDec) math.LegacyDec
```

## Parameters

The minting module contains the following parameters:

| Key                 | Type            | Example                      |
|---------------------|-----------------|------------------------------|
| MintDenom           | string          | "stake"                      |
| InflationsPerYear   | string (dec)    | "24625000000000000"          |
| BlocksPerYear       | string (uint64) | "6311520"                    |


## Events

The minting module emits the following events:

### BeginBlocker

| Type | Attribute Key     | Attribute Value    |
|------|-------------------|--------------------|
| mint | amount            | {amount}           |
