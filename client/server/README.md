# EVMStaking Module

## GetEVMStakingParams

URL: [GET] /evmstaking/params

### Response Example
```json
{
  "code": 200,
  "msg": {
    "params": {
      "max_withdrawal_per_block": 32,
      "max_sweep_per_block": 128,
      "min_partial_withdrawal_amount": "10000000",
      "ubi_withdraw_address": "0xcccccc0000000000000000000000000000000002"
    }
  },
  "error": ""
}
```

## GetWithdrawalQueue

URL: [GET] /evmstaking/withdrawal_queue

### Query Params
| Name                   | Type    | Example                      | Required |
|------------------------|---------|------------------------------|----------|
| pagination.key         | string  | FPoybu9dO+FCSV562u9keKVgUwur |          |
| pagination.offset      | string  | 0                            |          |
| pagination.limit       | string  | 10                           |          |
| pagination.count_total | boolean | true                         |          |
| pagination.reverse     | boolean | true                         |          |

### Response

- withdrawals: The list of withdrawals.
  - creation_height: The height of the withdrawal creation.
  - execution_address: The address to get the withdrawal.
  - amount: The amount of the withdrawal.
  - withdrawal_type: The type of the withdrawal.
    - 0: Unstake
    - 1: Reward
    - 2: UBI
  - validator_address: The address of the associated validator.
- pagination: The pagination info.
  - next_key: The key to query the next page.
  - total: The total number of withdrawals.

## GetRewardWithdrawalQueue

URL: [GET] /evmstaking/reward_withdrawal_queue

### Query Params
| Name                   | Type    | Example                      | Required |
|------------------------|---------|------------------------------|----------|
| pagination.key         | string  | FPoybu9dO+FCSV562u9keKVgUwur |          |
| pagination.offset      | string  | 0                            |          |
| pagination.limit       | string  | 10                           |          |
| pagination.count_total | boolean | true                         |          |
| pagination.reverse     | boolean | true                         |          |

### Response

- withdrawals: The list of withdrawals.
  - creation_height: The height of the withdrawal creation.
  - execution_address: The address to get the withdrawal.
  - amount: The amount of the withdrawal.
  - withdrawal_type: The type of the withdrawal.
    - 0: Unstake
    - 1: Reward
    - 2: UBI
  - validator_address: The address of the associated validator.
- pagination: The pagination info.
  - next_key: The key to query the next page.
  - total: The total number of withdrawals.

## GetOperatorAddress

URL: [GET] /evmstaking/delegators/{delegator_address}/operator_address

### Path Params
| Name              | Type   | Example                                      | Required |
|-------------------|--------|----------------------------------------------|----------|
| delegator_address | string | 0x00a842dbd3d11176b4868dd753a552b8919d5a63   | ✔        |

## GetWithdrawAddress

URL: [GET] /evmstaking/delegators/{delegator_address}/withdraw_address

### Path Params
| Name              | Type   | Example                                      | Required |
|-------------------|--------|----------------------------------------------|----------|
| delegator_address | string | 0x00a842dbd3d11176b4868dd753a552b8919d5a63   | ✔        |

## GetRewardAddress

URL: [GET] /evmstaking/delegators/{delegator_address}/reward_address

### Path Params
| Name              | Type   | Example                                      | Required |
|-------------------|--------|----------------------------------------------|----------|
| delegator_address | string | 0x00a842dbd3d11176b4868dd753a552b8919d5a63   | ✔        |

# Staking Module

## GetStakingParams

URL: [GET] /staking/params

### Response Example
```json
{
  "code": 200,
  "msg": {
    "params": {
      "unbonding_time": "10000000000",
      "max_validators": 32,
      "max_entries": 14,
      "historical_entries": 10000,
      "bond_denom": "stake",
      "min_commission_rate": "0.000000000000000000",
      "min_delegation": "1024",
      "periods": [
        {
          "duration": "0",
          "rewards_multiplier": "1.000000000000000000"
        },
        {
          "period_type": 1,
          "duration": "60000000000",
          "rewards_multiplier": "1.051000000000000000"
        },
        {
          "period_type": 2,
          "duration": "120000000000",
          "rewards_multiplier": "1.160000000000000000"
        },
        {
          "period_type": 3,
          "duration": "180000000000",
          "rewards_multiplier": "1.340000000000000000"
        }
      ],
      "token_types": [
        {
          "rewards_multiplier": "0.500000000000000000"
        },
        {
          "token_type": 1,
          "rewards_multiplier": "1.000000000000000000"
        }
      ]
    }
  },
  "error": ""
}
```

## GetStakingPool

URL: [GET] /staking/pool

### Response Example
```json
{
  "code": 200,
  "msg": {
    "pool": {
      "not_bonded_tokens": "76461600000000",
      "bonded_tokens": "80110471008000000"
    }
  },
  "error": ""
}
```

## GetHistoricalInfoByHeight

URL: [GET] /staking/historical_info/{height}

### Path Params
| Name   | Type    | Example | Required |
|--------|---------|---------|----------|
| height | integer | 1       | ✔        |

## GetTotalDelegatorsCount

URL: [GET] /staking/total_delegators_count

### Response Example
```json
{
    "code": 200,
    "msg": {
        "total": "100"
    },
    "error": ""
}
```

## GetValidators

URL: [GET] /staking/validators

### Query Params
| Name                   | Type    | Example                      | Required |
|------------------------|---------|------------------------------|----------|
| pagination.key         | string  | FPoybu9dO+FCSV562u9keKVgUwur |          |
| pagination.offset      | string  | 0                            |          |
| pagination.limit       | string  | 10                           |          |
| pagination.count_total | boolean | true                         |          |
| pagination.reverse     | boolean | true                         |          |
| status                 | string  | 3                            |          |

### Response Example
```json
{
  "code": 200,
  "msg": {
    "validators": [
      {
        "operator_address": "0x00a842dbd3d11176b4868dd753a552b8919d5a63",
        "consensus_pubkey": {
          "type": "tendermint/PubKeySecp256k1",
          "value": "AvLo+lkg0UWozoI+pJzv1a7upt+HaMxZCdWgRxvZ8Cb1"
        },
        "status": 3,
        "tokens": "1000000",
        "delegator_shares": "1000000.000000000000000000",
        "description": {
          "moniker": "0x0FC41199CE588948861A8DA86D725A5A073AE91A"
        },
        "unbonding_time": "1970-01-01T00:00:00Z",
        "commission": {
          "commission_rates": {
            "rate": "0.070000000000000000",
            "max_rate": "0.100000000000000000",
            "max_change_rate": "0.010000000000000000"
          },
          "update_time": "2024-04-16T11:04:40.60280319Z"
        },
        "min_self_delegation": "1024000000000",
        "support_token_type": 1,
        "rewards_tokens": "1000000.000000000000000000",
        "delegator_rewards_shares": "1000000.000000000000000000"
      },
      {
        "operator_address": "0x13665369a8ad5163f0c023839323b5d015925de1",
        "consensus_pubkey": {
          "type": "tendermint/PubKeySecp256k1",
          "value": "AlZ/RRQCXVnTPjZWqxXdaZ0X0ZvJvJRDHu1R4UpgwzXl"
        },
        "status": 3,
        "tokens": "1000000",
        "delegator_shares": "1000000.000000000000000000",
        "description": {
          "moniker": "0x99C28AE30CBEFEFF75E91C66692FE0BD9279B861"
        },
        "unbonding_time": "1970-01-01T00:00:00Z",
        "commission": {
          "commission_rates": {
            "rate": "0.070000000000000000",
            "max_rate": "0.100000000000000000",
            "max_change_rate": "0.010000000000000000"
          },
          "update_time": "2024-04-16T11:04:40.60280319Z"
        },
        "min_self_delegation": "1024000000000",
        "support_token_type": 1,
        "rewards_tokens": "1000000.000000000000000000",
        "delegator_rewards_shares": "1000000.000000000000000000"
      },
      {
        "operator_address": "0x87f3cc50c84005f7130d37b849f6a71e05a8bf1f",
        "consensus_pubkey": {
          "type": "tendermint/PubKeySecp256k1",
          "value": "AzUWGMooEM92H8RCIOqXbjtbeur+2rOzgP9T/umnf1eA"
        },
        "status": 3,
        "tokens": "1000000",
        "delegator_shares": "1000000.000000000000000000",
        "description": {
          "moniker": "0x9DFC26A7662106EEEC5E87B20CBB690CFCE73A05"
        },
        "unbonding_time": "1970-01-01T00:00:00Z",
        "commission": {
          "commission_rates": {
            "rate": "0.070000000000000000",
            "max_rate": "0.100000000000000000",
            "max_change_rate": "0.010000000000000000"
          },
          "update_time": "2024-04-16T11:04:40.60280319Z"
        },
        "min_self_delegation": "1024000000000",
        "support_token_type": 1,
        "rewards_tokens": "1000000.000000000000000000",
        "delegator_rewards_shares": "1000000.000000000000000000"
      },
      {
        "operator_address": "0xc47c28f925179089b6b7b1b336ac1f943b240066",
        "consensus_pubkey": {
          "type": "tendermint/PubKeySecp256k1",
          "value": "AqCVQjEtIzN9q5sMgSgl4dDD27vx6wa528lp9rjNKZE/"
        },
        "status": 3,
        "tokens": "1000000",
        "delegator_shares": "1000000.000000000000000000",
        "description": {
          "moniker": "0x768A39103B552E7AE56635DD4E9B55922AAFC504"
        },
        "unbonding_time": "1970-01-01T00:00:00Z",
        "commission": {
          "commission_rates": {
            "rate": "0.070000000000000000",
            "max_rate": "0.100000000000000000",
            "max_change_rate": "0.010000000000000000"
          },
          "update_time": "2024-04-16T11:04:40.60280319Z"
        },
        "min_self_delegation": "1024000000000",
        "support_token_type": 1,
        "rewards_tokens": "1000000.000000000000000000",
        "delegator_rewards_shares": "1000000.000000000000000000"
      }
    ],
    "pagination": {
      "total": "4"
    }
  },
  "error": ""
}
```

## GetValidator

URL: [GET] /staking/validators/{validator_addr}

### Path Params
| Name           | Type   | Example                                      | Required |
|----------------|--------|----------------------------------------------|----------|
| validator_addr | string | 0x00a842dbd3d11176b4868dd753a552b8919d5a63   | ✔        |

### Response Example
```json
{
  "code": 200,
  "msg": {
    "validator": {
      "operator_address": "0x00a842dbd3d11176b4868dd753a552b8919d5a63",
      "consensus_pubkey": {
        "type": "tendermint/PubKeySecp256k1",
        "value": "AvLo+lkg0UWozoI+pJzv1a7upt+HaMxZCdWgRxvZ8Cb1"
      },
      "status": 3,
      "tokens": "1000000",
      "delegator_shares": "1000000.000000000000000000",
      "description": {
        "moniker": "0x0FC41199CE588948861A8DA86D725A5A073AE91A"
      },
      "unbonding_time": "1970-01-01T00:00:00Z",
      "commission": {
        "commission_rates": {
          "rate": "0.070000000000000000",
          "max_rate": "0.100000000000000000",
          "max_change_rate": "0.010000000000000000"
        },
        "update_time": "2024-04-16T11:04:40.60280319Z"
      },
      "min_self_delegation": "1024000000000",
      "support_token_type": 1,
      "rewards_tokens": "1000000.000000000000000000",
      "delegator_rewards_shares": "1000000.000000000000000000"
    }
  },
  "error": ""
}
```

## GetValidatorDelegations

URL: [GET] /staking/validators/{validator_addr}/delegations

### Path Params
| Name           | Type   | Example                                      | Required |
|----------------|--------|----------------------------------------------|----------|
| validator_addr | string | 0x00a842dbd3d11176b4868dd753a552b8919d5a63   | ✔        |

### Query Params
| Name                   | Type    | Example                      | Required |
|------------------------|---------|------------------------------|----------|
| pagination.key         | string  | FPoybu9dO+FCSV562u9keKVgUwur |          |
| pagination.offset      | string  | 0                            |          |
| pagination.limit       | string  | 10                           |          |
| pagination.count_total | boolean | true                         |          |
| pagination.reverse     | boolean | true                         |          |

### Response Example
```json
{
  "code": 200,
  "msg": {
    "delegation_responses": [
      {
        "delegation": {
          "delegator_address": "0x00a842dbd3d11176b4868dd753a552b8919d5a63",
          "validator_address": "0x00a842dbd3d11176b4868dd753a552b8919d5a63",
          "shares": "1000000.000000000000000000",
          "rewards_shares": "1000000.000000000000000000"
        },
        "balance": {
          "denom": "stake",
          "amount": "1000000"
        }
      }
    ],
    "pagination": {
      "total": "1"
    }
  },
  "error": ""
}
```

## GetValidatorDelegation

URL: [GET] /staking/validators/{validator_addr}/delegations/{delegator_addr}

### Path Params
| Name           | Type   | Example                                      | Required |
|----------------|--------|----------------------------------------------|----------|
| validator_addr | string | 0x00a842dbd3d11176b4868dd753a552b8919d5a63   | ✔        |
| delegator_addr | string | 0x00a842dbd3d11176b4868dd753a552b8919d5a63   | ✔        |

### Response Example
```json
{
  "code": 200,
  "msg": {
    "delegation_response": {
      "delegation": {
        "delegator_address": "0x00a842dbd3d11176b4868dd753a552b8919d5a63",
        "validator_address": "0x00a842dbd3d11176b4868dd753a552b8919d5a63",
        "shares": "1000000.000000000000000000",
        "rewards_shares": "1000000.000000000000000000"
      },
      "balance": {
        "denom": "stake",
        "amount": "1000000"
      }
    }
  },
  "error": ""
}
```

## GetValidatorTotalDelegationsCount

URL: [GET] /staking/validators/{validator_addr}/total_delegations_count

### Path Params
| Name           | Type   | Example                                      | Required |
|----------------|--------|----------------------------------------------|----------|
| validator_addr | string | 0x00a842dbd3d11176b4868dd753a552b8919d5a63   | ✔        |

### Response Example
```json
{
    "code": 200,
    "msg": {
        "total": "11"
    },
    "error": ""
}
```

## GetValidatorPeriodDelegations

URL: [GET] /staking/validators/{validator_addr}/delegators/{delegator_addr}/period_delegations

### Path Params
| Name           | Type   | Example                                      | Required |
|----------------|--------|----------------------------------------------|----------|
| validator_addr | string | 0x00a842dbd3d11176b4868dd753a552b8919d5a63   | ✔        |
| delegator_addr | string | 0x00a842dbd3d11176b4868dd753a552b8919d5a63   | ✔        |

### Query Params
| Name                   | Type    | Example                      | Required |
|------------------------|---------|------------------------------|----------|
| pagination.key         | string  | FPoybu9dO+FCSV562u9keKVgUwur |          |
| pagination.offset      | string  | 0                            |          |
| pagination.limit       | string  | 10                           |          |
| pagination.count_total | boolean | true                         |          |
| pagination.reverse     | boolean | true                         |          |

### Response Example
```json
{
  "code": 200,
  "msg": {
    "period_delegation_responses": [
      {
        "period_delegation": {
          "delegator_address": "0x00a842dbd3d11176b4868dd753a552b8919d5a63",
          "validator_address": "0x00a842dbd3d11176b4868dd753a552b8919d5a63",
          "period_delegation_id": "0",
          "shares": "1025000000000.000000000000000000",
          "rewards_shares": "1025000000000.000000000000000000",
          "end_time": "2024-10-23T08:48:00.313756096Z"
        },
        "balance": {
          "denom": "stake",
          "amount": "1025000000000.000000000000000000"
        }
      }
    ],
    "pagination": {
      "next_key": "FN80AOspT3VnYmzjLPoXTD2DEC9B",
      "total": "66"
    }
  },
  "error": ""
}
```

## GetValidatorPeriodDelegation

URL: [GET] /staking/validators/{validator_addr}/delegators/{delegator_addr}/period_delegations/{period_delegation_id}

### Path Params
| Name                 | Type   | Example                                      | Required |
|----------------------|--------|----------------------------------------------|----------|
| validator_addr       | string | 0x00a842dbd3d11176b4868dd753a552b8919d5a63   | ✔        |
| delegator_addr       | string | 0x00a842dbd3d11176b4868dd753a552b8919d5a63   | ✔        |
| period_delegation_id | string | 0                                            | ✔        |

### Response Example
```json
{
  "code": 200,
  "msg": {
    "period_delegation_response": {
      "period_delegation": {
        "delegator_address": "0x00a842dbd3d11176b4868dd753a552b8919d5a63",
        "validator_address": "0x00a842dbd3d11176b4868dd753a552b8919d5a63",
        "period_delegation_id": "0",
        "shares": "1025000000000.000000000000000000",
        "rewards_shares": "1025000000000.000000000000000000",
        "end_time": "2024-10-23T08:48:00.313756096Z"
      },
      "balance": {
        "denom": "stake",
        "amount": "1025000000000.000000000000000000"
      }
    }
  },
  "error": ""
}
```

## GetValidatorUnbondingDelegations

URL: [GET] /staking/validators/{validator_addr}/unbonding_delegations

### Path Params
| Name           | Type   | Example                                      | Required |
|----------------|--------|----------------------------------------------|----------|
| validator_addr | string | 0x00a842dbd3d11176b4868dd753a552b8919d5a63   | ✔        |

### Query Params
| Name                   | Type       | Example                                    | Required |
|------------------------|------------|--------------------------------------------|----------|
| pagination.key         | string     | FPoybu9dO+FCSV562u9keKVgUwur               |          |
| pagination.offset      | string     | 0                                          |          |
| pagination.limit       | string     | 10                                         |          |
| pagination.count_total | boolean    | true                                       |          |
| pagination.reverse     | boolean    | true                                       |          |

### Response Example
```json
{
  "code": 200,
  "msg": {
    "unbonding_responses": [
      {
        "delegator_address": "0x00a842dbd3d11176b4868dd753a552b8919d5a63",
        "validator_address": "0x00a842dbd3d11176b4868dd753a552b8919d5a63",
        "entries": [
          {
            "creation_height": "525632",
            "completion_time": "2024-11-27T08:26:13.41935718Z",
            "initial_balance": "1024000000000",
            "balance": "1024000000000",
            "unbonding_id": "53"
          }
        ]
      }
    ],
    "pagination": {
      "total": "1"
    }
  },
  "error": ""
}
```

## GetValidatorUnbondingDelegation

URL: [GET] /staking/validators/{validator_addr}/delegations/{delegator_addr}/unbonding_delegation

### Path Params
| Name           | Type   | Example                                      | Required |
|----------------|--------|----------------------------------------------|----------|
| validator_addr | string | 0x00a842dbd3d11176b4868dd753a552b8919d5a63   | ✔        |
| delegator_addr | string | 0x00a842dbd3d11176b4868dd753a552b8919d5a63   | ✔        |

## GetDelegatorValidators

URL: [GET] /staking/delegators/{delegator_addr}/validators

### Path Params
| Name           | Type   | Example                                      | Required |
|----------------|--------|----------------------------------------------|----------|
| delegator_addr | string | 0x00a842dbd3d11176b4868dd753a552b8919d5a63   | ✔        |

### Query Params
| Name                   | Type    | Example                                    | Required |
|------------------------|---------|--------------------------------------------|----------|
| pagination.key         | string  | FPoybu9dO+FCSV562u9keKVgUwur               |          |
| pagination.offset      | string  | 0                                          |          |
| pagination.limit       | array   | 10                                         |          |
| pagination.count_total | boolean | true                                       |          |
| pagination.reverse     | boolean | true                                       |          |

### Response Example
```json
{
  "code": 200,
  "msg": {
    "validators": [
      {
        "operator_address": "0x00a842dbd3d11176b4868dd753a552b8919d5a63",
        "consensus_pubkey": {
          "type": "tendermint/PubKeySecp256k1",
          "value": "AvLo+lkg0UWozoI+pJzv1a7upt+HaMxZCdWgRxvZ8Cb1"
        },
        "status": 3,
        "tokens": "1000000",
        "delegator_shares": "1000000.000000000000000000",
        "description": {
          "moniker": "0x0FC41199CE588948861A8DA86D725A5A073AE91A"
        },
        "unbonding_time": "1970-01-01T00:00:00Z",
        "commission": {
          "commission_rates": {
            "rate": "0.070000000000000000",
            "max_rate": "0.100000000000000000",
            "max_change_rate": "0.010000000000000000"
          },
          "update_time": "2024-04-16T11:04:40.60280319Z"
        },
        "min_self_delegation": "1024000000000",
        "support_token_type": 1,
        "rewards_tokens": "1000000.000000000000000000",
        "delegator_rewards_shares": "1000000.000000000000000000"
      }
    ],
    "pagination": {
      "total": "1"
    }
  },
  "error": ""
}
```

## GetDelegatorValidator

URL: [GET] /staking/delegators/{delegator_addr}/validators/{validator_addr}

### Path Params
| Name           | Type   | Example                                      | Required |
|----------------|--------|----------------------------------------------|----------|
| delegator_addr | string | 0x00a842dbd3d11176b4868dd753a552b8919d5a63   | ✔        |
| validator_addr | string | 0x00a842dbd3d11176b4868dd753a552b8919d5a63   | ✔        |

### Response Example
```json
{
  "code": 200,
  "msg": {
    "validator": {
      "operator_address": "0x00a842dbd3d11176b4868dd753a552b8919d5a63",
      "consensus_pubkey": {
        "type": "tendermint/PubKeySecp256k1",
        "value": "A3ljgjHsCTGAxbqOtMjkQ66DudEruUyCrXFgYkYI4geh"
      },
      "jailed": true,
      "status": 1,
      "tokens": "5722800000000",
      "delegator_shares": "6024000000000.000000000000000000",
      "description": {
        "moniker": "validator-04-create"
      },
      "unbonding_height": "4859",
      "unbonding_time": "2024-10-22T12:32:13.342436482Z",
      "commission": {
        "commission_rates": {
          "rate": "0.100000000000000000",
          "max_rate": "0.500000000000000000",
          "max_change_rate": "0.010000000000000000"
        },
        "update_time": "2024-10-22T12:16:08.070117612Z"
      },
      "min_self_delegation": "1024",
      "unbonding_ids": [
        "147"
      ],
      "support_token_type": 1,
      "rewards_tokens": "5772412800000.000000000000000000",
      "delegator_rewards_shares": "6076224000000.000000000000000000"
    }
  },
  "error": ""
}
```

## GetDelegatorStakedToken

URL: [GET] /staking/delegators/{delegator_addr}/staked_token

### Path Params
| Name           | Type   | Example                                      | Required |
|----------------|--------|----------------------------------------------|----------|
| delegator_addr | string | 0x00a842dbd3d11176b4868dd753a552b8919d5a63   | ✔        |

### Query Params
| Name                   | Type    | Example                                    | Required |
|------------------------|---------|--------------------------------------------|----------|
| pagination.key         | string  | FPoybu9dO+FCSV562u9keKVgUwur               |          |
| pagination.offset      | string  | 0                                          |          |
| pagination.limit       | array   | 10                                         |          |
| pagination.count_total | boolean | true                                       |          |
| pagination.reverse     | boolean | true                                       |          |

### Response Example
```json
{
  "code": 200,
  "msg": {
    "delegation_staked_token": [
      {
        "validator_operator_address": "0x88355450c9003d1d43773bd72b837e818693a781",
        "staked_token": "2048000000000.000000000000000000"
      }
    ],
    "pagination": {
      "total": "1"
    }
  },
  "error": ""
}
```

## GetDelegatorTotalStakedToken

URL: [GET] /staking/delegators/{delegator_addr}/total_staked_token

### Path Params
| Name           | Type   | Example                                      | Required |
|----------------|--------|----------------------------------------------|----------|
| delegator_addr | string | 0x00a842dbd3d11176b4868dd753a552b8919d5a63   | ✔        |

### Response Example
```json
{
  "code": 200,
  "msg": {
    "staked_token": "2048000000000.000000000000000000"
  },
  "error": ""
}
```

## GetDelegatorRewardsToken

URL: [GET] /staking/delegators/{delegator_addr}/rewards_token

### Path Params
| Name           | Type   | Example                                      | Required |
|----------------|--------|----------------------------------------------|----------|
| delegator_addr | string | 0x00a842dbd3d11176b4868dd753a552b8919d5a63   | ✔        |

### Query Params
| Name                   | Type    | Example                                    | Required |
|------------------------|---------|--------------------------------------------|----------|
| pagination.key         | string  | FPoybu9dO+FCSV562u9keKVgUwur               |          |
| pagination.offset      | string  | 0                                          |          |
| pagination.limit       | array   | 10                                         |          |
| pagination.count_total | boolean | true                                       |          |
| pagination.reverse     | boolean | true                                       |          |

### Response Example
```json
{
    "code": 200,
    "msg": {
        "delegation_rewards_token": [
            {
                "validator_operator_address": "0xc5c0beeac8b37ed52f6a675ee2154d926a88e3ec",
                "rewards_token": "58709824999896.500000000000000000"
            },
            {
                "validator_operator_address": "0xcd29b70ff04c0aa386f7b3453df0e5ed3d4f67bb",
                "rewards_token": "136740151396798.000000000000000000"
            },
            {
                "validator_operator_address": "0xcd5faabca5bea3c5fc5e2371c7b397604720c2c2",
                "rewards_token": "58708000000000.000000000000000000"
            },
            {
                "validator_operator_address": "0xdb8e606ad7c02f37e43d10a10126791dc94b0434",
                "rewards_token": "119408704000000.000000000000000000"
            }
        ],
        "pagination": {
            "total": "4"
        }
    },
    "error": ""
}
```

## GetDelegatorTotalRewardsToken

URL: [GET] /staking/delegators/{delegator_addr}/total_rewards_token

### Path Params
| Name           | Type   | Example                                      | Required |
|----------------|--------|----------------------------------------------|----------|
| delegator_addr | string | 0x00a842dbd3d11176b4868dd753a552b8919d5a63   | ✔        |

### Response Example
```json
{
    "code": 200,
    "msg": {
        "rewards_token": "373566680396694.500000000000000000"
    },
    "error": ""
}
```

## GetDelegatorDelegations

URL: [GET] /staking/delegations/{delegator_addr}

### Path Params
| Name           | Type   | Example                                      | Required |
|----------------|--------|----------------------------------------------|----------|
| delegator_addr | string | 0x00a842dbd3d11176b4868dd753a552b8919d5a63   | ✔        |

### Query Params
| Name                   | Type    | Example                                    | Required |
|------------------------|---------|--------------------------------------------|----------|
| pagination.key         | string  | FPoybu9dO+FCSV562u9keKVgUwur               |          |
| pagination.offset      | string  | 0                                          |          |
| pagination.limit       | string  | 10                                         |          |
| pagination.count_total | boolean | true                                       |          |
| pagination.reverse     | boolean | true                                       |          |

### Response Example
```json
{
  "code": 200,
  "msg": {
    "delegation_responses": [
      {
        "delegation": {
          "delegator_address": "0x00a842dbd3d11176b4868dd753a552b8919d5a63",
          "validator_address": "0x00a842dbd3d11176b4868dd753a552b8919d5a63",
          "shares": "1000000.000000000000000000",
          "rewards_shares": "1000000.000000000000000000"
        },
        "balance": {
          "denom": "stake",
          "amount": "1000000"
        }
      }
    ],
    "pagination": {
      "total": "1"
    }
  },
  "error": ""
}
```

## GetDelegatorRedelegations

URL: [GET] /staking/delegators/{delegator_addr}/redelegations

### Path Params
| Name           | Type   | Example                                      | Required |
|----------------|--------|----------------------------------------------|----------|
| delegator_addr | string | 0x00a842dbd3d11176b4868dd753a552b8919d5a63   | ✔        |

### Query Params
| Name                   | Type    | Example                                    | Required |
|------------------------|---------|--------------------------------------------|----------|
| pagination.key         | string  | FPoybu9dO+FCSV562u9keKVgUwur               |          |
| pagination.offset      | string  | 0                                          |          |
| pagination.limit       | array   | 10                                         |          |
| pagination.count_total | boolean | true                                       |          |
| pagination.reverse     | boolean | true                                       |          |
| src_validator_addr     | string  | 0x87f3cc50c84005f7130d37b849f6a71e05a8bf1f |          |
| dst_validator_addr     | string  | 0xc47c28f925179089b6b7b1b336ac1f943b240066 |          |

## GetDelegatorUnbondingDelegations

URL: [GET] /staking/delegators/{delegator_addr}/unbonding_delegations

### Path Params
| Name           | Type   | Example                                      | Required |
|----------------|--------|----------------------------------------------|----------|
| delegator_addr | string | 0x00a842dbd3d11176b4868dd753a552b8919d5a63   | ✔        |

### Query Params
| Name                   | Type    | Example                                    | Required |
|------------------------|---------|--------------------------------------------|----------|
| pagination.key         | string  | FPoybu9dO+FCSV562u9keKVgUwur               |          |
| pagination.offset      | string  | 0                                          |          |
| pagination.limit       | array   | 10                                         |          |
| pagination.count_total | boolean | true                                       |          |
| pagination.reverse     | boolean | true                                       |          |

# Auth Module

## GetAuthParams

URL: [GET] /auth/params

### Response Example
```json
{
  "code": 200,
  "msg": {
    "params": {
      "max_memo_characters": "256",
      "tx_sig_limit": "7",
      "tx_size_cost_per_byte": "10",
      "sig_verify_cost_ed25519": "590",
      "sig_verify_cost_secp256k1": "1000"
    }
  },
  "error": ""
}
```

# Bank Module

## GetBankParams

URL: [GET] /bank/params

### Response Example
```json
{
  "code": 200,
  "msg": {
    "params": {}
  },
  "error": ""
}
```

## GetSupplyByDenom

URL: [GET] /bank/supply/by_denom

### Query Params
| Name  | Type   | Example | Required |
|-------|--------|---------|----------|
| denom | string | stake   | ✔        |

### Response Example
```json
{
  "code": 200,
  "msg": {
    "amount": {
      "denom": "stake",
      "amount": "79748037587244392"
    }
  },
  "error": ""
}
```

## GetBalancesByAddressDenom

URL: [GET] /bank/balances/{address}/by_denom

### Path Params
| Name    | Type   | Example                                      | Required |
|---------|--------|----------------------------------------------|----------|
| address | string | 0x00a842dbd3d11176b4868dd753a552b8919d5a63   | ✔        |

### Query Params
| Name  | Type   | Example | Required |
|-------|--------|---------|----------|
| denom | string | stake   |          |

### Response Example
```json
{
  "code": 200,
  "msg": {
    "balance": {
      "denom": "stake",
      "amount": "50301746"
    }
  },
  "error": ""
}
```

## GetSpendableBalancesByAddressDenom

URL: [GET] /bank/spendable_balances/{address}/by_denom

### Path Params
| Name    | Type   | Example                                      | Required |
|---------|--------|----------------------------------------------|----------|
| address | string | 0x00a842dbd3d11176b4868dd753a552b8919d5a63   | ✔        |

### Query Params
| Name  | Type   | Example | Required |
|-------|--------|---------|----------|
| denom | string | stake   |          |

### Response Example
```json
{
  "code": 200,
  "msg": {
    "balance": {
      "denom": "stake",
      "amount": "50301746"
    }
  },
  "error": ""
}
```

# Distribution Module

## GetDistributionParams

URL: [GET] /distribution/params

### Response Example
```json
{
  "code": 200,
  "msg": {
    "params": {
      "ubi": "0.065000000000000000",
      "base_proposer_reward": "0.000000000000000000",
      "bonus_proposer_reward": "0.000000000000000000",
      "withdraw_addr_enabled": true
    }
  },
  "error": ""
}
```

## GetValidatorCommissionByValidatorAddress

URL: [GET] /distribution/validators/{validator_address}/commission

### Path Params
| Name              | Type   | Example                                      | Required |
|-------------------|--------|----------------------------------------------|----------|
| validator_address | string | 0x00a842dbd3d11176b4868dd753a552b8919d5a63   | ✔        |

### Response Example
```json
{
  "code": 200,
  "msg": {
    "commission": {
      "commission": [
        {
          "denom": "stake",
          "amount": "11129810.318708930942711499"
        }
      ]
    }
  },
  "error": ""
}
```

## GetValidatorOutstandingRewardsByValidatorAddress

URL: [GET] /distribution/validators/{validator_address}/outstanding_rewards

### Path Params
| Name             | Type   | Example                                      | Required |
|------------------|--------|----------------------------------------------|----------|
| validator_address| string | 0x00a842dbd3d11176b4868dd753a552b8919d5a63   | ✔        |

### Response Example
```json
{
  "code": 200,
  "msg": {
    "rewards": {
      "rewards": [
        {
          "denom": "stake",
          "amount": "110356851.470660927219114823"
        }
      ]
    }
  },
  "error": ""
}
```

## GetValidatorSlashesByValidatorAddress

URL: [GET] /distribution/validators/{validator_address}/slashes

### Path Params
| Name             | Type   | Example                                      | Required |
|------------------|--------|----------------------------------------------|----------|
| validator_address| string | 0x00a842dbd3d11176b4868dd753a552b8919d5a63   | ✔        |

### Query Params
| Name                   | Type    | Example                                    | Required |
|------------------------|---------|--------------------------------------------|----------|
| pagination.key         | string  | FPoybu9dO+FCSV562u9keKVgUwur               |          |
| pagination.offset      | string  | 0                                          |          |
| pagination.limit       | string  | 10                                         |          |
| pagination.count_total | boolean | true                                       |          |
| pagination.reverse     | boolean | true                                       |          |
| starting_height        | string  | 1                                          |          |
| ending_height          | string  | 100                                        |          |

### Response Example
```json
{
  "code": 200,
  "msg": {
    "slashes": [],
    "pagination": {}
  },
  "error": ""
}
```

## GetDistributionValidatorByValidatorAddress

URL: [GET] /distribution/validators/{validator_address}

### Path Params
| Name             | Type   | Example                                      | Required |
|------------------|--------|----------------------------------------------|----------|
| validator_address| string | 0x00a842dbd3d11176b4868dd753a552b8919d5a63   | ✔        |

### Response Example
```json
{
  "code": 200,
  "msg": {
    "operator_address": "0x00a842dbd3d11176b4868dd753a552b8919d5a63",
    "self_bond_rewards": [
      {
        "denom": "stake",
        "amount": "82232804.757899250000000000"
      }
    ],
    "commission": [
      {
        "denom": "stake",
        "amount": "11129810.318708930942711499"
      }
    ]
  },
  "error": ""
}
```

## GetDistributionValidatorsByDelegatorAddress

URL: [GET] /distribution/delegators/{delegator_address}/validators

### Path Params
| Name              | Type   | Example                                      | Required |
|-------------------|--------|----------------------------------------------|----------|
| delegator_address | string | story1f5zuqhmwy39cv64g6laeeg264ydz06txlfqtg3 | ✔        |

### Response Example
```json
{
  "code": 200,
  "msg": {
    "validators": [
      "0x00a842dbd3d11176b4868dd753a552b8919d5a63"
    ]
  },
  "error": ""
}
```

## GetDelegatorRewardsByDelegatorAddress

URL: [GET] /distribution/delegators/{delegator_address}/rewards

### Path Params
| Name              | Type   | Example                                      | Required |
|-------------------|--------|----------------------------------------------|----------|
| delegator_address | string | 0x00a842dbd3d11176b4868dd753a552b8919d5a63   | ✔        |

### Response Example
```json
{
  "code": 200,
  "msg": {
    "rewards": [
      {
        "validator_address": "0x00a842dbd3d11176b4868dd753a552b8919d5a63",
        "reward": [
          {
            "denom": "stake",
            "amount": "519497.220124590080000000"
          }
        ]
      }
    ],
    "total": [
      {
        "denom": "stake",
        "amount": "519497.220124590080000000"
      }
    ]
  },
  "error": ""
}
```

## GetDelegatorRewardsByDelegatorAddressValidatorAddress

URL: [GET] /distribution/delegators/{delegator_address}/rewards/{validator_address}

### Path Params
| Name              | Type   | Example                                      | Required |
|-------------------|--------|----------------------------------------------|----------|
| delegator_address | string | 0x00a842dbd3d11176b4868dd753a552b8919d5a63   | ✔        |
| validator_address | string | 0x00a842dbd3d11176b4868dd753a552b8919d5a63   | ✔        |

### Response Example
```json
{
  "code": 200,
  "msg": {
    "rewards": [
      {
        "denom": "stake",
        "amount": "16994236.393259084134400000"
      }
    ]
  },
  "error": ""
}
```

# Mint Module

## GetMintParams

URL: [GET] /mint/params

### Response Example
```json
{
  "code": 200,
  "msg": {
    "params": {
      "mint_denom": "stake",
      "inflations_per_year": "24625000000000000.000000000000000000",
      "blocks_per_year": "10368000"
    }
  },
  "error": ""
}
```

# Slashing Module

## GetSlashingParams

URL: [GET] /slashing/params

### Response Example
```json
{
  "code": 200,
  "msg": {
    "params": {
      "signed_blocks_window": "200",
      "min_signed_per_window": "0.050000000000000000",
      "downtime_jail_duration": "60000000000",
      "slash_fraction_double_sign": "0.050000000000000000",
      "slash_fraction_downtime": "0.050000000000000000"
    }
  },
  "error": ""
}
```

## GetSigningInfo

URL: [GET] /slashing/signing_infos/{pubkey}

### Path Params
| Name   | Type   | Example                                                                 | Required |
|--------|--------|-------------------------------------------------------------------------|----------|
| pubkey | string | 0x03351618ca2810cf761fc44220ea976e3b5b7aeafedab3b380ff53fee9a77f5780    | ✔        |

# Upgrade Module

## GetAppliedPlan

URL: [GET] /upgrade/applied_plan/{name}

### Path Params
| Name | Type   | Example   | Required |
|------|--------|-----------|----------|
| name | string | upgrade-1 |    ✔     |

## GetAuthority

URL: [GET] /upgrade/authority

## GetCurrentPlan

URL: [GET] /upgrade/current_plan

## GetModuleVersions

URL: [GET] /upgrade/module_versions

### Query Params
| Name        | Type   | Example | Required |
|-------------|--------|---------|----------|
| module_name | string | staking |    ✔     |
