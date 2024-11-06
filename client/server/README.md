## GetValidators

URL: [GET] /staking/validators

### Query Params
| Name | Type | Example | Require |
| --- | --- | --- | --- |
| status | string |  | ✔ |
| pagination.key | string |  | ✔ |
| pagination.offset | string | 0 | ✔ |
| pagination.limit | string | 10 | ✔ |
| pagination.count_total | string | true | ✔ |
| pagination.reverse | string | true | ✔ |

### Response Example
```json
{
  "code": 200,
  "msg": {
    "validators": [
      {
        "operator_address": "storyvaloper1l43wgqtf825k4qwx4lt6xremu0jjptv7w6q0ak",
        "consensus_pubkey": {
          "type": "tendermint/PubKeySecp256k1",
          "value": "Avihb1CqLW0pX2YJ6ru0/N7sr5MWnk5JCy/U5ezQ9axo"
        },
        "jailed": true,
        "status": 1,
        "tokens": "4750000000000",
        "delegator_shares": "5000000000000.000000000000000000",
        "description": {
          "moniker": "validator-test104-21"
        },
        "unbonding_height": "4351",
        "unbonding_time": "2024-10-21T14:12:46.362852882Z",
        "commission": {
          "commission_rates": {
            "rate": "0.100000000000000000",
            "max_rate": "0.500000000000000000",
            "max_change_rate": "0.010000000000000000"
          },
          "update_time": "2024-10-21T13:37:04.758645464Z"
        },
        "min_self_delegation": "1024",
        "unbonding_ids": [
          "245"
        ],
        "rewards_tokens": "2375000000000.000000000000000000",
        "delegator_rewards_shares": "2500000000000.000000000000000000"
      },
      {
        "operator_address": "storyvaloper17mmqx2r9d5zq8z9jcajxpc64tkhh5tuxzl5rdl",
        "consensus_pubkey": {
          "type": "tendermint/PubKeySecp256k1",
          "value": "AxnWLFO8LIeRpbnrKFVY2P+vPzG6bI5tpkivg/F+0XM6"
        },
        "jailed": true,
        "status": 1,
        "tokens": "4750000000000",
        "delegator_shares": "5000000000000.000000000000000000",
        "description": {
          "moniker": "validator-test104-24"
        },
        "unbonding_height": "4351",
        "unbonding_time": "2024-10-21T14:12:46.362852882Z",
        "commission": {
          "commission_rates": {
            "rate": "0.100000000000000000",
            "max_rate": "0.500000000000000000",
            "max_change_rate": "0.010000000000000000"
          },
          "update_time": "2024-10-21T13:38:30.541469298Z"
        },
        "min_self_delegation": "1024",
        "unbonding_ids": [
          "244"
        ],
        "rewards_tokens": "2375000000000.000000000000000000",
        "delegator_rewards_shares": "2500000000000.000000000000000000"
      },
      {
        "operator_address": "storyvaloper17myx8sae0dxyd50je73j5fx9e7gemffefe5hey",
        "consensus_pubkey": {
          "type": "tendermint/PubKeySecp256k1",
          "value": "Ak1HAUyvr5iPJztkNrISe7G06p244pnfzZXJqhFMUO7a"
        },
        "jailed": true,
        "status": 1,
        "tokens": "4750000000000",
        "delegator_shares": "5000000000000.000000000000000000",
        "description": {
          "moniker": "validator-test104-16"
        },
        "unbonding_height": "4344",
        "unbonding_time": "2024-10-21T14:12:24.142758147Z",
        "commission": {
          "commission_rates": {
            "rate": "0.100000000000000000",
            "max_rate": "0.500000000000000000",
            "max_change_rate": "0.010000000000000000"
          },
          "update_time": "2024-10-21T13:34:51.559498051Z"
        },
        "min_self_delegation": "1024",
        "unbonding_ids": [
          "181",
          "243"
        ],
        "rewards_tokens": "2375000000000.000000000000000000",
        "delegator_rewards_shares": "2500000000000.000000000000000000"
      },
      {
        "operator_address": "storyvaloper17z6fh2fqdqpa9am3a8m7c024r954j7hwgzdu0z",
        "consensus_pubkey": {
          "type": "tendermint/PubKeySecp256k1",
          "value": "AgFUP1Y/JMO0GR4KvOxNRmYzzU1VThmVkEORPpvsaGk0"
        },
        "jailed": true,
        "status": 1,
        "tokens": "1024000000000",
        "delegator_shares": "1024000000000.000000000000000000",
        "description": {
          "moniker": "test-0x8F51B4fA7C60016A8e1Db1caE312E0dF63F21E00"
        },
        "unbonding_height": "1166",
        "unbonding_time": "2024-10-21T11:24:12.073004276Z",
        "commission": {
          "commission_rates": {
            "rate": "0.100000000000000000",
            "max_rate": "0.500000000000000000",
            "max_change_rate": "0.010000000000000000"
          },
          "update_time": "2024-10-21T11:22:55.156217405Z"
        },
        "min_self_delegation": "1024",
        "unbonding_ids": [
          "29"
        ],
        "support_token_type": 1,
        "rewards_tokens": "1076224000000.000000000000000000",
        "delegator_rewards_shares": "1076224000000.000000000000000000"
      },
      {
        "operator_address": "storyvaloper1amf8h5xfrh2hz0zweke42t6m86zrl7ycyzgd60",
        "consensus_pubkey": {
          "type": "tendermint/PubKeySecp256k1",
          "value": "AooAFKj/AvQhjSgyylAeebrb8RgsbtMZxN3fefzjy2hf"
        },
        "jailed": true,
        "status": 1,
        "tokens": "4750000000000",
        "delegator_shares": "5000000000000.000000000000000000",
        "description": {
          "moniker": "validator-test104-13"
        },
        "unbonding_height": "4318",
        "unbonding_time": "2024-10-21T14:11:01.350669293Z",
        "commission": {
          "commission_rates": {
            "rate": "0.100000000000000000",
            "max_rate": "0.500000000000000000",
            "max_change_rate": "0.010000000000000000"
          },
          "update_time": "2024-10-21T13:33:25.868660453Z"
        },
        "min_self_delegation": "1024",
        "unbonding_ids": [
          "180",
          "242"
        ],
        "rewards_tokens": "2375000000000.000000000000000000",
        "delegator_rewards_shares": "2500000000000.000000000000000000"
      },
      {
        "operator_address": "storyvaloper1awyv4u5pknxpjayqs50sm5p6yq8p4xy7awavph",
        "consensus_pubkey": {
          "type": "tendermint/PubKeySecp256k1",
          "value": "A32c3gy91iNyQ2BRSWz7gIusPWG0JI3Qp3D6k85p/W9J"
        },
        "jailed": true,
        "status": 1,
        "tokens": "972800000000",
        "delegator_shares": "1024000000000.000000000000000000",
        "description": {
          "moniker": "test-0x8D83d84315B98ACbA32A364Ae841635d650CD6Ab"
        },
        "unbonding_height": "1598",
        "unbonding_time": "2024-10-21T11:47:02.079884724Z",
        "commission": {
          "commission_rates": {
            "rate": "0.100000000000000000",
            "max_rate": "0.500000000000000000",
            "max_change_rate": "0.010000000000000000"
          },
          "update_time": "2024-10-21T11:30:57.629535295Z"
        },
        "min_self_delegation": "1024",
        "unbonding_ids": [
          "60"
        ],
        "rewards_tokens": "486400000000.000000000000000000",
        "delegator_rewards_shares": "512000000000.000000000000000000"
      },
      {
        "operator_address": "storyvaloper1af4dyn4q27a460lxqnynp3zp8lgahsvl0gt6v8",
        "consensus_pubkey": {
          "type": "tendermint/PubKeySecp256k1",
          "value": "Ag9YT4iMqthseWCDZPU81ZNyAbIu6kRZKk3WafTXisyP"
        },
        "status": 3,
        "tokens": "6024000000000",
        "delegator_shares": "6024000000000.000000000000000000",
        "description": {
          "moniker": "validator-01-create"
        },
        "unbonding_time": "1970-01-01T00:00:00Z",
        "commission": {
          "commission_rates": {
            "rate": "0.100000000000000000",
            "max_rate": "0.500000000000000000",
            "max_change_rate": "0.010000000000000000"
          },
          "update_time": "2024-10-21T14:51:28.064815617Z"
        },
        "min_self_delegation": "1024",
        "support_token_type": 1,
        "rewards_tokens": "6372160000000.000000000000000000",
        "delegator_rewards_shares": "6372160000000.000000000000000000"
      },
      {
        "operator_address": "storyvaloper1u6skd9rmmlrzdl8ttwxywyzvapwxe7y67xcrvl",
        "consensus_pubkey": {
          "type": "tendermint/PubKeySecp256k1",
          "value": "AldIN1AzKGNcp0oLkuodo9s1huqMXj17B4wyjHE2bh3C"
        },
        "jailed": true,
        "status": 1,
        "tokens": "4750000000000",
        "delegator_shares": "5000000000000.000000000000000000",
        "description": {
          "moniker": "validator-test104-10"
        },
        "unbonding_height": "4301",
        "unbonding_time": "2024-10-21T14:10:07.22524099Z",
        "commission": {
          "commission_rates": {
            "rate": "0.100000000000000000",
            "max_rate": "0.500000000000000000",
            "max_change_rate": "0.010000000000000000"
          },
          "update_time": "2024-10-21T13:32:12.738388066Z"
        },
        "min_self_delegation": "1024",
        "unbonding_ids": [
          "182",
          "241"
        ],
        "rewards_tokens": "2375000000000.000000000000000000",
        "delegator_rewards_shares": "2500000000000.000000000000000000"
      },
      {
        "operator_address": "storyvaloper1u60eyhv44w578au4lv982wjkcuhhhyss2ukqj9",
        "consensus_pubkey": {
          "type": "tendermint/PubKeySecp256k1",
          "value": "A9cvfpWD3qASAFaCzavpbSd7E+IRKKTDZkYBvxY5gRk7"
        },
        "status": 3,
        "tokens": "1005000001000000",
        "delegator_shares": "1005000001000000.000000000000000000",
        "description": {
          "moniker": "0xE69F925D95ABA9E3F795FB0A753A56C72F7B9210"
        },
        "unbonding_time": "1970-01-01T00:00:00Z",
        "commission": {
          "commission_rates": {
            "rate": "0.000000000000000000",
            "max_rate": "0.000000000000000000",
            "max_change_rate": "0.000000000000000000"
          },
          "update_time": "2024-04-16T11:04:40.60280319Z"
        },
        "min_self_delegation": "1024",
        "support_token_type": 1,
        "rewards_tokens": "1005000001000000.000000000000000000",
        "delegator_rewards_shares": "1005000001000000.000000000000000000"
      },
      {
        "operator_address": "storyvaloper1ue8pas8awavc2duhv5hrqt4d2stk7d9d7czzef",
        "consensus_pubkey": {
          "type": "tendermint/PubKeySecp256k1",
          "value": "AzElmMH8Ce73/O/az2G+YJcSpHPWWi4oP6ISN8Zd75mT"
        },
        "jailed": true,
        "status": 1,
        "tokens": "4750000000000",
        "delegator_shares": "5000000000000.000000000000000000",
        "description": {
          "moniker": "validator-test104-29"
        },
        "unbonding_height": "4281",
        "unbonding_time": "2024-10-21T14:09:03.681537469Z",
        "commission": {
          "commission_rates": {
            "rate": "0.100000000000000000",
            "max_rate": "0.500000000000000000",
            "max_change_rate": "0.010000000000000000"
          },
          "update_time": "2024-10-21T13:40:43.850630572Z"
        },
        "min_self_delegation": "1024",
        "unbonding_ids": [
          "240"
        ],
        "rewards_tokens": "2375000000000.000000000000000000",
        "delegator_rewards_shares": "2500000000000.000000000000000000"
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


## GetValidatorByValidatorAddress

URL: [GET] /staking/validators/{validator_addr}

### Path Params
| Name | Type | Example | Require |
| --- | --- | --- | --- |
| validator_addr | string | storyvaloper17973uudyv484cgmy4sd4kjwdgrh86vfwhh0nlw | ✔ |

### Response Example
```json
{
  "code": 200,
  "msg": {
    "validator": {
      "operator_address": "storyvaloper17973uudyv484cgmy4sd4kjwdgrh86vfwhh0nlw",
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


## GetDelegatorByDelegatorAddress

URL: [GET] /staking/delegators/{delegator_addr}

### Path Params
| Name | Type | Example | Require |
| --- | --- | --- | --- |
| delegator_addr | string | story1f5zuqhmwy39cv64g6laeeg264ydz06txlfqtg3 | ✔ |

### Response Example
```json
{
  "code": 200,
  "msg": {
    "delegator_addr": "story1l3ewm920n34gj2cam7lcfangyutn84574cpnl2",
    "withdraw_address": "0x6E5e0eFC2961ed81e663E53395C51b4154855e77",
    "reward_address": "0x6E5e0eFC2961ed81e663E53395C51b4154855e77",
    "operator_address": "0x6E5e0eFC2961ed81e663E53395C51b4154855e77"
  },
  "error": ""
}
```


## GetValidatorDelegationsByValidatorAddress

URL: [GET] /staking/validators/{validator_addr}/delegations

### Path Params
| Name | Type | Example | Require |
| --- | --- | --- | --- |
| validator_addr | string | storyvaloper17973uudyv484cgmy4sd4kjwdgrh86vfwhh0nlw | ✔ |

### Response Example
```json
{
  "code": 200,
  "msg": {
    "delegation_responses": [
      {
        "delegation": {
          "delegator_address": "story1f5zuqhmwy39cv64g6laeeg264ydz06txlfqtg3",
          "validator_address": "storyvaloper17973uudyv484cgmy4sd4kjwdgrh86vfwhh0nlw",
          "shares": "1024000000000.000000000000000000",
          "rewards_shares": "1076224000000.000000000000000000"
        },
        "balance": {
          "denom": "stake",
          "amount": "972800000000"
        }
      },
      {
        "delegation": {
          "delegator_address": "story17973uudyv484cgmy4sd4kjwdgrh86vfwecmj59",
          "validator_address": "storyvaloper17973uudyv484cgmy4sd4kjwdgrh86vfwhh0nlw",
          "shares": "5000000000000.000000000000000000",
          "rewards_shares": "5000000000000.000000000000000000"
        },
        "balance": {
          "denom": "stake",
          "amount": "4750000000000"
        }
      }
    ],
    "pagination": {
      "total": "2"
    }
  },
  "error": ""
}
```


## GetDelegationsByDelegatorAddress

URL: [GET] /staking/delegations/{delegator_addr}

### Path Params
| Name | Type | Example | Require |
| --- | --- | --- | --- |
| delegator_addr | string | story1f5zuqhmwy39cv64g6laeeg264ydz06txlfqtg3 | ✔ |

### Response Example
```json
{
  "code": 200,
  "msg": {
    "delegation_responses": [
      {
        "delegation": {
          "delegator_address": "story1f5zuqhmwy39cv64g6laeeg264ydz06txlfqtg3",
          "validator_address": "storyvaloper1rk277ak9enq422sevnf2yjjc6mhdv44wl8u30p",
          "shares": "1024000000000.000000000000000000",
          "rewards_shares": "1187840000000.000000000000000000"
        },
        "balance": {
          "denom": "stake",
          "amount": "1024000000000"
        }
      },
      {
        "delegation": {
          "delegator_address": "story1f5zuqhmwy39cv64g6laeeg264ydz06txlfqtg3",
          "validator_address": "storyvaloper1wk04lv2e2egfdt25e94xvwszny2ww8znancs5z",
          "shares": "4096000000000.000000000000000000",
          "rewards_shares": "2048000000000.000000000000000000"
        },
        "balance": {
          "denom": "stake",
          "amount": "4096000000000"
        }
      },
      {
        "delegation": {
          "delegator_address": "story1f5zuqhmwy39cv64g6laeeg264ydz06txlfqtg3",
          "validator_address": "storyvaloper137yg6cc4gpvhdx5wt0q3g5pcjeljuqvs4eu8y2",
          "shares": "1024000000000.000000000000000000",
          "rewards_shares": "1076224000000.000000000000000000"
        },
        "balance": {
          "denom": "stake",
          "amount": "972800000000"
        }
      },
      {
        "delegation": {
          "delegator_address": "story1f5zuqhmwy39cv64g6laeeg264ydz06txlfqtg3",
          "validator_address": "storyvaloper1hcx6v7887r7r7498ge3kxqxrmkt3ekpk7jj5c2",
          "shares": "2048000000000.000000000000000000",
          "rewards_shares": "2375680000000.000000000000000000"
        },
        "balance": {
          "denom": "stake",
          "amount": "2048000000000"
        }
      },
      {
        "delegation": {
          "delegator_address": "story1f5zuqhmwy39cv64g6laeeg264ydz06txlfqtg3",
          "validator_address": "storyvaloper1mwgwgn6x5waumhg4tgnuhkezhamnaq4385mrgv",
          "shares": "2048000000000.000000000000000000",
          "rewards_shares": "1024000000000.000000000000000000"
        },
        "balance": {
          "denom": "stake",
          "amount": "1945600000000"
        }
      },
      {
        "delegation": {
          "delegator_address": "story1f5zuqhmwy39cv64g6laeeg264ydz06txlfqtg3",
          "validator_address": "storyvaloper17973uudyv484cgmy4sd4kjwdgrh86vfwhh0nlw",
          "shares": "1024000000000.000000000000000000",
          "rewards_shares": "1076224000000.000000000000000000"
        },
        "balance": {
          "denom": "stake",
          "amount": "972800000000"
        }
      }
    ],
    "pagination": {
      "total": "6"
    }
  },
  "error": ""
}
```


## GetRedelegationsByDelegatorAddress

URL: [GET] /staking/delegators/{delegator_addr}/redelegations

### Path Params
| Name | Type | Example | Require |
| --- | --- | --- | --- |
| delegator_addr | string | story1f5zuqhmwy39cv64g6laeeg264ydz06txlfqtg3 | ✔ |

### Query Params
| Name | Type | Example | Require |
| --- | --- | --- | --- |
| status | string |  | ✔ |
| pagination.key | string |  | ✔ |
| pagination.offset | string | 0 | ✔ |
| pagination.limit | string | 10 | ✔ |
| pagination.count_total | string | true | ✔ |
| pagination.reverse | string | true | ✔ |



## GetValidatorsByDelegatorAddress

URL: [GET] /staking/delegators/{delegator_addr}/validators

### Path Params
| Name | Type | Example | Require |
| --- | --- | --- | --- |
| delegator_addr | string | story1f5zuqhmwy39cv64g6laeeg264ydz06txlfqtg3 | ✔ |

### Response Example
```json
{
  "code": 200,
  "msg": {
    "validators": [
      {
        "operator_address": "storyvaloper1rk277ak9enq422sevnf2yjjc6mhdv44wl8u30p",
        "consensus_pubkey": {
          "type": "tendermint/PubKeySecp256k1",
          "value": "A/TbagXrHvqfxf+9XVDqU51QAFKAgzLEqyAFHKN57jmR"
        },
        "jailed": true,
        "status": 1,
        "tokens": "1024000000000",
        "delegator_shares": "1024000000000.000000000000000000",
        "description": {
          "moniker": "test-0x4Be8b027f918627Ec1C2D017dB041d40494fd80A"
        },
        "unbonding_height": "3680",
        "unbonding_time": "2024-10-22T11:29:42.978993465Z",
        "commission": {
          "commission_rates": {
            "rate": "0.100000000000000000",
            "max_rate": "0.500000000000000000",
            "max_change_rate": "0.010000000000000000"
          },
          "update_time": "2024-10-22T11:28:51.693143157Z"
        },
        "min_self_delegation": "1024",
        "unbonding_ids": [
          "77"
        ],
        "support_token_type": 1,
        "rewards_tokens": "1187840000000.000000000000000000",
        "delegator_rewards_shares": "1187840000000.000000000000000000"
      },
      {
        "operator_address": "storyvaloper1wk04lv2e2egfdt25e94xvwszny2ww8znancs5z",
        "consensus_pubkey": {
          "type": "tendermint/PubKeySecp256k1",
          "value": "A2FY1BIrKwcrSPCOUWESOEJBt52nkKAb6sfML8A5oDRF"
        },
        "jailed": true,
        "status": 1,
        "tokens": "4096000000000",
        "delegator_shares": "4096000000000.000000000000000000",
        "description": {
          "moniker": "validator-test"
        },
        "unbonding_height": "5576",
        "unbonding_time": "2024-10-22T13:10:07.737947097Z",
        "commission": {
          "commission_rates": {
            "rate": "0.100000000000000000",
            "max_rate": "0.500000000000000000",
            "max_change_rate": "0.010000000000000000"
          },
          "update_time": "2024-10-22T13:08:16.216165198Z"
        },
        "min_self_delegation": "1024",
        "unbonding_ids": [
          "189"
        ],
        "rewards_tokens": "2048000000000.000000000000000000",
        "delegator_rewards_shares": "2048000000000.000000000000000000"
      },
      {
        "operator_address": "storyvaloper137yg6cc4gpvhdx5wt0q3g5pcjeljuqvs4eu8y2",
        "consensus_pubkey": {
          "type": "tendermint/PubKeySecp256k1",
          "value": "AjaQELCH8iw4p7DfBCKVn7ZEK+R0PHUEh36z3TMte7jO"
        },
        "jailed": true,
        "status": 1,
        "tokens": "5722800000000",
        "delegator_shares": "6024000000000.000000000000000000",
        "description": {
          "moniker": "validator-01-create"
        },
        "unbonding_height": "4853",
        "unbonding_time": "2024-10-22T12:31:54.211161149Z",
        "commission": {
          "commission_rates": {
            "rate": "0.100000000000000000",
            "max_rate": "0.500000000000000000",
            "max_change_rate": "0.010000000000000000"
          },
          "update_time": "2024-10-22T12:15:48.838926766Z"
        },
        "min_self_delegation": "1024",
        "unbonding_ids": [
          "146"
        ],
        "support_token_type": 1,
        "rewards_tokens": "5772412800000.000000000000000000",
        "delegator_rewards_shares": "6076224000000.000000000000000000"
      },
      {
        "operator_address": "storyvaloper1hcx6v7887r7r7498ge3kxqxrmkt3ekpk7jj5c2",
        "consensus_pubkey": {
          "type": "tendermint/PubKeySecp256k1",
          "value": "AwjVZXRcfWsQhFgYngZyTd8P601uYsMLlaVNP48XgG48"
        },
        "jailed": true,
        "status": 1,
        "tokens": "2048000000000",
        "delegator_shares": "2048000000000.000000000000000000",
        "description": {
          "moniker": "test-0x426CaA619dF2B01831A08C0af8075482be20D83C"
        },
        "unbonding_height": "5037",
        "unbonding_time": "2024-10-22T12:41:37.967378773Z",
        "commission": {
          "commission_rates": {
            "rate": "0.100000000000000000",
            "max_rate": "0.500000000000000000",
            "max_change_rate": "0.010000000000000000"
          },
          "update_time": "2024-10-22T12:35:51.661691725Z"
        },
        "min_self_delegation": "1024",
        "unbonding_ids": [
          "158"
        ],
        "support_token_type": 1,
        "rewards_tokens": "2375680000000.000000000000000000",
        "delegator_rewards_shares": "2375680000000.000000000000000000"
      },
      {
        "operator_address": "storyvaloper1mwgwgn6x5waumhg4tgnuhkezhamnaq4385mrgv",
        "consensus_pubkey": {
          "type": "tendermint/PubKeySecp256k1",
          "value": "A3egvqcBqnapCQ0Vj4CbwC1J2Vbkxajf8RQOV54RnJz+"
        },
        "jailed": true,
        "status": 1,
        "tokens": "1945600000000",
        "delegator_shares": "2048000000000.000000000000000000",
        "description": {
          "moniker": "test-0x83Aa3657d70df740Efac57457DB7c24c92f3cFbf"
        },
        "unbonding_height": "7919",
        "unbonding_time": "2024-10-22T15:19:51.227049925Z",
        "commission": {
          "commission_rates": {
            "rate": "0.100000000000000000",
            "max_rate": "0.500000000000000000",
            "max_change_rate": "0.010000000000000000"
          },
          "update_time": "2024-10-22T14:16:18.419084723Z"
        },
        "min_self_delegation": "1024",
        "unbonding_ids": [
          "246",
          "254"
        ],
        "rewards_tokens": "972800000000.000000000000000000",
        "delegator_rewards_shares": "1024000000000.000000000000000000"
      },
      {
        "operator_address": "storyvaloper17973uudyv484cgmy4sd4kjwdgrh86vfwhh0nlw",
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
    ],
    "pagination": {
      "total": "6"
    }
  },
  "error": ""
}
```


## GetValidatorsByDelegatorAddressValidatorAddress

URL: [GET] /staking/delegators/{delegator_addr}/validators/{validator_addr}

### Path Params
| Name | Type | Example | Require |
| --- | --- | --- | --- |
| delegator_addr | string | story1f5zuqhmwy39cv64g6laeeg264ydz06txlfqtg3 | ✔ |
| validator_addr | string | storyvaloper17973uudyv484cgmy4sd4kjwdgrh86vfwhh0nlw | ✔ |

### Response Example
```json
{
  "code": 200,
  "msg": {
    "validator": {
      "operator_address": "storyvaloper17973uudyv484cgmy4sd4kjwdgrh86vfwhh0nlw",
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


## GetUnbondingDelegationsByDelegatorAddress

URL: [GET] /staking/delegators/{delegator_addr}/unbonding_delegations

### Path Params
| Name | Type | Example | Require |
| --- | --- | --- | --- |
| delegator_addr | string | story1f5zuqhmwy39cv64g6laeeg264ydz06txlfqtg3 | ✔ |

### Query Params
| Name | Type | Example | Require |
| --- | --- | --- | --- |
| status | string |  | ✔ |
| pagination.key | string |  | ✔ |
| pagination.offset | string | 0 | ✔ |
| pagination.limit | string | 10 | ✔ |
| pagination.count_total | string | true | ✔ |
| pagination.reverse | string | true | ✔ |



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
| Name | Type | Example | Require |
| --- | --- | --- | --- |
| height | integer |  | ✔ |

### Response Example
```json
{
  "code": 200,
  "msg": {
    "hist": {
      "header": {
        "version": {},
        "chain_id": "odyssey-devnet-1",
        "height": "6666",
        "time": "2024-10-22T14:09:17.546670477Z",
        "last_block_id": {
          "part_set_header": {}
        },
        "next_validators_hash": "bHIf/wtJa7OBLolGYkRxgiAr72c708itMxqPhAwC0O4=",
        "app_hash": "UDRmZ+Tg6ZMW0ZIuErfYXE1VS+QnWbAj8BVm9p07f3c=",
        "proposer_address": "Y/nE5kiyRMMZVjDDr0YiArKC4w4="
      },
      "valset": [
        {
          "operator_address": "storyvaloper1w69rjypm25h84etxxhw5ax64jg42l3gyhr0j42",
          "consensus_pubkey": {
            "type": "tendermint/PubKeySecp256k1",
            "value": "A9mdRUGE+sv2oD6jfrNvalDGmELqOtQgOKjVU3vRWyWU"
          },
          "status": 3,
          "tokens": "10104447001000000",
          "delegator_shares": "10104447001000000.000000000000000000",
          "description": {
            "moniker": "0x768A39103B552E7AE56635DD4E9B55922AAFC504"
          },
          "unbonding_time": "1970-01-01T00:00:00Z",
          "commission": {
            "commission_rates": {
              "rate": "0.000000000000000000",
              "max_rate": "0.000000000000000000",
              "max_change_rate": "0.000000000000000000"
            },
            "update_time": "2024-04-16T11:04:40.60280319Z"
          },
          "min_self_delegation": "1024",
          "rewards_tokens": "5052223500500000.000000000000000000",
          "delegator_rewards_shares": "5052223500500000.000000000000000000"
        },
        {
          "operator_address": "storyvaloper1q9vn6g9fj96rw2anyklt56v96rf20ap5z7zha7",
          "consensus_pubkey": {
            "type": "tendermint/PubKeySecp256k1",
            "value": "A3e5PxDT1PEstFxhGtoUShU32oIKCnnXtBp3IBi9bM/2"
          },
          "status": 3,
          "tokens": "10000000001000000",
          "delegator_shares": "10000000001000000.000000000000000000",
          "description": {
            "moniker": "0x01593D20A99174372BB325BEBA6985D0D2A7F434"
          },
          "unbonding_time": "1970-01-01T00:00:00Z",
          "commission": {
            "commission_rates": {
              "rate": "0.000000000000000000",
              "max_rate": "0.000000000000000000",
              "max_change_rate": "0.000000000000000000"
            },
            "update_time": "2024-04-16T11:04:40.60280319Z"
          },
          "min_self_delegation": "1024",
          "support_token_type": 1,
          "rewards_tokens": "10000000001000000.000000000000000000",
          "delegator_rewards_shares": "10000000001000000.000000000000000000"
        },
        {
          "operator_address": "storyvaloper1plzprxwwtzy53ps63k5x6uj6tgrn46g6c8w4xv",
          "consensus_pubkey": {
            "type": "tendermint/PubKeySecp256k1",
            "value": "AqBVHHkyOfiie29Wrez6hMvC644kbZfPgXA1jFEs7Uwq"
          },
          "status": 3,
          "tokens": "10000000001000000",
          "delegator_shares": "10000000001000000.000000000000000000",
          "description": {
            "moniker": "0x0FC41199CE588948861A8DA86D725A5A073AE91A"
          },
          "unbonding_time": "1970-01-01T00:00:00Z",
          "commission": {
            "commission_rates": {
              "rate": "0.000000000000000000",
              "max_rate": "0.000000000000000000",
              "max_change_rate": "0.000000000000000000"
            },
            "update_time": "2024-04-16T11:04:40.60280319Z"
          },
          "min_self_delegation": "1024",
          "rewards_tokens": "5000000000500000.000000000000000000",
          "delegator_rewards_shares": "5000000000500000.000000000000000000"
        },
        {
          "operator_address": "storyvaloper1v0uufejgkfzvxx2kxrp6733zq2eg9ccwgx5emd",
          "consensus_pubkey": {
            "type": "tendermint/PubKeySecp256k1",
            "value": "AtJpE37ydCffhJT0NNkFm/8CEB64cYu5xBVs/btKqjWh"
          },
          "status": 3,
          "tokens": "10000000001000000",
          "delegator_shares": "10000000001000000.000000000000000000",
          "description": {
            "moniker": "0x63F9C4E648B244C3195630C3AF462202B282E30E"
          },
          "unbonding_time": "1970-01-01T00:00:00Z",
          "commission": {
            "commission_rates": {
              "rate": "0.000000000000000000",
              "max_rate": "0.000000000000000000",
              "max_change_rate": "0.000000000000000000"
            },
            "update_time": "2024-04-16T11:04:40.60280319Z"
          },
          "min_self_delegation": "1024",
          "support_token_type": 1,
          "rewards_tokens": "10000000001000000.000000000000000000",
          "delegator_rewards_shares": "10000000001000000.000000000000000000"
        },
        {
          "operator_address": "storyvaloper106cyk87evdmf3z2rmkj5m9aut5xcchpwsrlrs7",
          "consensus_pubkey": {
            "type": "tendermint/PubKeySecp256k1",
            "value": "A+wUiPhFG58EShf9w6v6BDmdey3rsst5S17jR14fezbz"
          },
          "status": 3,
          "tokens": "10000000001000000",
          "delegator_shares": "10000000001000000.000000000000000000",
          "description": {
            "moniker": "0x7EB04B1FD96376988943DDA54D97BC5D0D8C5C2E"
          },
          "unbonding_time": "1970-01-01T00:00:00Z",
          "commission": {
            "commission_rates": {
              "rate": "0.000000000000000000",
              "max_rate": "0.000000000000000000",
              "max_change_rate": "0.000000000000000000"
            },
            "update_time": "2024-04-16T11:04:40.60280319Z"
          },
          "min_self_delegation": "1024",
          "support_token_type": 1,
          "rewards_tokens": "10000000001000000.000000000000000000",
          "delegator_rewards_shares": "10000000001000000.000000000000000000"
        },
        {
          "operator_address": "storyvaloper1n8pg4ccvhml07a0fr3nxjtlqhkf8nwrpc6a2t5",
          "consensus_pubkey": {
            "type": "tendermint/PubKeySecp256k1",
            "value": "A/9SMxZTnh3Rq96Eygg9MfB6g82euMXhjT5nMWrhLlyf"
          },
          "status": 3,
          "tokens": "10000000001000000",
          "delegator_shares": "10000000001000000.000000000000000000",
          "description": {
            "moniker": "0x99C28AE30CBEFEFF75E91C66692FE0BD9279B861"
          },
          "unbonding_time": "1970-01-01T00:00:00Z",
          "commission": {
            "commission_rates": {
              "rate": "0.000000000000000000",
              "max_rate": "0.000000000000000000",
              "max_change_rate": "0.000000000000000000"
            },
            "update_time": "2024-04-16T11:04:40.60280319Z"
          },
          "min_self_delegation": "1024",
          "rewards_tokens": "5000000000500000.000000000000000000",
          "delegator_rewards_shares": "5000000000500000.000000000000000000"
        },
        {
          "operator_address": "storyvaloper1nh7zdfmxyyrwamz7s7eqewmfpn7wwws9gyue2t",
          "consensus_pubkey": {
            "type": "tendermint/PubKeySecp256k1",
            "value": "A6KRGirXFYsv5oVQz8d8YIl0Nj23bXo2jLHui72y12Bi"
          },
          "status": 3,
          "tokens": "10000000001000000",
          "delegator_shares": "10000000001000000.000000000000000000",
          "description": {
            "moniker": "0x9DFC26A7662106EEEC5E87B20CBB690CFCE73A05"
          },
          "unbonding_time": "1970-01-01T00:00:00Z",
          "commission": {
            "commission_rates": {
              "rate": "0.000000000000000000",
              "max_rate": "0.000000000000000000",
              "max_change_rate": "0.000000000000000000"
            },
            "update_time": "2024-04-16T11:04:40.60280319Z"
          },
          "min_self_delegation": "1024",
          "rewards_tokens": "5000000000500000.000000000000000000",
          "delegator_rewards_shares": "5000000000500000.000000000000000000"
        },
        {
          "operator_address": "storyvaloper1u60eyhv44w578au4lv982wjkcuhhhyss2ukqj9",
          "consensus_pubkey": {
            "type": "tendermint/PubKeySecp256k1",
            "value": "A9cvfpWD3qASAFaCzavpbSd7E+IRKKTDZkYBvxY5gRk7"
          },
          "status": 3,
          "tokens": "10000000001000000",
          "delegator_shares": "10000000001000000.000000000000000000",
          "description": {
            "moniker": "0xE69F925D95ABA9E3F795FB0A753A56C72F7B9210"
          },
          "unbonding_time": "1970-01-01T00:00:00Z",
          "commission": {
            "commission_rates": {
              "rate": "0.000000000000000000",
              "max_rate": "0.000000000000000000",
              "max_change_rate": "0.000000000000000000"
            },
            "update_time": "2024-04-16T11:04:40.60280319Z"
          },
          "min_self_delegation": "1024",
          "support_token_type": 1,
          "rewards_tokens": "10000000001000000.000000000000000000",
          "delegator_rewards_shares": "10000000001000000.000000000000000000"
        },
        {
          "operator_address": "storyvaloper1rg0vphgl9nzy8njgkclxkr9wu9tvh4raly2eex",
          "consensus_pubkey": {
            "type": "tendermint/PubKeySecp256k1",
            "value": "AgSPN1XgoqtMqu0JJb5SC76FXvKO70Y+KXh3xjfqMalw"
          },
          "status": 3,
          "tokens": "6024000000000",
          "delegator_shares": "6024000000000.000000000000000000",
          "description": {
            "moniker": "validator-test"
          },
          "unbonding_time": "1970-01-01T00:00:00Z",
          "commission": {
            "commission_rates": {
              "rate": "0.100000000000000000",
              "max_rate": "0.500000000000000000",
              "max_change_rate": "0.010000000000000000"
            },
            "update_time": "2024-10-22T13:54:22.256319155Z"
          },
          "min_self_delegation": "1024",
          "support_token_type": 1,
          "rewards_tokens": "6024000000000.000000000000000000",
          "delegator_rewards_shares": "6024000000000.000000000000000000"
        }
      ]
    }
  },
  "error": ""
}
```


## GetDelegationByValidatorAddressDelegatorAddress

URL: [GET] /staking/validators/{validator_addr}/delegations/{delegator_addr}

### Path Params
| Name | Type | Example | Require |
| --- | --- | --- | --- |
| validator_addr | string | storyvaloper17973uudyv484cgmy4sd4kjwdgrh86vfwhh0nlw | ✔ |
| delegator_addr | string | story1f5zuqhmwy39cv64g6laeeg264ydz06txlfqtg3 | ✔ |

### Response Example
```json
{
  "code": 200,
  "msg": {
    "delegation_response": {
      "delegation": {
        "delegator_address": "story1f5zuqhmwy39cv64g6laeeg264ydz06txlfqtg3",
        "validator_address": "storyvaloper17973uudyv484cgmy4sd4kjwdgrh86vfwhh0nlw",
        "shares": "1024000000000.000000000000000000",
        "rewards_shares": "1076224000000.000000000000000000"
      },
      "balance": {
        "denom": "stake",
        "amount": "972800000000"
      }
    }
  },
  "error": ""
}
```


## GetValidatorUnbondingDelegations

URL: [GET] /staking/validators/{validator_addr}/unbonding_delegations

### Path Params
| Name | Type | Example | Require |
| --- | --- | --- | --- |
| validator_addr | string | storyvaloper17973uudyv484cgmy4sd4kjwdgrh86vfwhh0nlw | ✔ |

### Query Params
| Name | Type | Example | Require |
| --- | --- | --- | --- |
| status | string |  | ✔ |
| pagination.key | string |  | ✔ |
| pagination.offset | string | 0 | ✔ |
| pagination.limit | string | 10 | ✔ |
| pagination.count_total | string | true | ✔ |
| pagination.reverse | string | true | ✔ |



## GetDelegatorUnbondingDelegation

URL: [GET] /staking/validators/{validator_addr}/delegations/{delegator_addr}/unbonding_delegation

### Path Params
| Name | Type | Example | Require |
| --- | --- | --- | --- |
| validator_addr | string | storyvaloper17973uudyv484cgmy4sd4kjwdgrh86vfwhh0nlw | ✔ |
| delegator_addr | string | story1f5zuqhmwy39cv64g6laeeg264ydz06txlfqtg3 | ✔ |



## GetPeriodDelegationsByDelegatorAddress

URL: [GET] /staking/validators/{validator_addr}/delegators/{delegator_addr}/period_delegations

### Path Params
| Name | Type | Example | Require |
| --- | --- | --- | --- |
| validator_addr | string | storyvaloper17973uudyv484cgmy4sd4kjwdgrh86vfwhh0nlw | ✔ |
| delegator_addr | string | story1f5zuqhmwy39cv64g6laeeg264ydz06txlfqtg3 | ✔ |

### Response Example
```json
{
  "code": 200,
  "msg": [
    {
      "delegator_address": "story1l3ewm920n34gj2cam7lcfangyutn84574cpnl2",
      "validator_address": "storyvaloper1l3ewm920n34gj2cam7lcfangyutn8457mh4j5p",
      "period_delegation_id": "0",
      "shares": "1025000000000.000000000000000000",
      "rewards_shares": "1025000000000.000000000000000000",
      "end_time": "2024-10-23T08:48:00.313756096Z"
    }
  ],
  "error": ""
}
```


## GetPeriodDelegationByDelegatorAddressAndID

URL: [GET] /staking/validators/{validator_addr}/delegators/{delegator_addr}/period_delegations/{period_delegation_id}

### Path Params
| Name | Type | Example | Require |
| --- | --- | --- | --- |
| validator_addr | string | storyvaloper17973uudyv484cgmy4sd4kjwdgrh86vfwhh0nlw | ✔ |
| delegator_addr | string | story1f5zuqhmwy39cv64g6laeeg264ydz06txlfqtg3 | ✔ |
| period_delegation_id | string | 1 | ✔ |

### Response Example
```json
{
  "code": 200,
  "msg": {
    "delegator_address": "story1l3ewm920n34gj2cam7lcfangyutn84574cpnl2",
    "validator_address": "storyvaloper1l3ewm920n34gj2cam7lcfangyutn8457mh4j5p",
    "period_delegation_id": "0",
    "shares": "1025000000000.000000000000000000",
    "rewards_shares": "1025000000000.000000000000000000",
    "end_time": "2024-10-23T08:48:00.313756096Z"
  },
  "error": ""
}
```


## GetAccounts

URL: [GET] /auth/accounts

### Query Params
| Name | Type | Example | Require |
| --- | --- | --- | --- |
| pagination.key | string |  | ✔ |
| pagination.offset | string | 0 | ✔ |
| pagination.limit | string | 10 | ✔ |
| pagination.count_total | string | true | ✔ |
| pagination.reverse | string | true | ✔ |

### Response Example
```json
{
  "code": 200,
  "msg": {
    "accounts": [
      {
        "type": "cosmos-sdk/BaseAccount",
        "value": {
          "address": "story1le74erycve0n78lg6ylunmkp886dwxl26qgz80",
          "account_number": "56"
        }
      },
      {
        "type": "cosmos-sdk/BaseAccount",
        "value": {
          "address": "story1lxp43exras9rkk6kpqxq5phgv5ttd60el5vfx9",
          "account_number": "20"
        }
      },
      {
        "type": "cosmos-sdk/BaseAccount",
        "value": {
          "address": "story176qh0f6jq3p5jzwykllhdzcj4twacjjeapvu9c",
          "account_number": "118"
        }
      },
      {
        "type": "cosmos-sdk/BaseAccount",
        "value": {
          "address": "story170yj2r7selm3dc0kkej8yfp0nwtgt456hgt5md",
          "account_number": "30"
        }
      },
      {
        "type": "cosmos-sdk/ModuleAccount",
        "value": {
          "address": "story17xpfvakm2amg962yls6f84z3kell8c5ljcetf8",
          "public_key": "",
          "account_number": 8,
          "sequence": 0,
          "name": "fee_collector",
          "permissions": null
        }
      },
      {
        "type": "cosmos-sdk/BaseAccount",
        "value": {
          "address": "story17973uudyv484cgmy4sd4kjwdgrh86vfwecmj59",
          "account_number": "52"
        }
      },
      {
        "type": "cosmos-sdk/BaseAccount",
        "value": {
          "address": "story17zuu6pujj4jfzzg35qxcrj4ee0jgd9xdftv6z2",
          "account_number": "70"
        }
      },
      {
        "type": "cosmos-sdk/BaseAccount",
        "value": {
          "address": "story17z6fh2fqdqpa9am3a8m7c024r954j7hwxdeayf",
          "account_number": "38"
        }
      },
      {
        "type": "cosmos-sdk/BaseAccount",
        "value": {
          "address": "story1a4zcqjltpwl3zfk72kx7m2pycflgcxatn44dy6",
          "account_number": "13"
        }
      },
      {
        "type": "cosmos-sdk/BaseAccount",
        "value": {
          "address": "story1anes54tdw9p7a35st7gfakqfffkqnm8d2n52l3",
          "account_number": "119"
        }
      }
    ],
    "pagination": {
      "next_key": "66ZhqK1PgBZTxEFTuO05iVtURCo=",
      "total": "120"
    }
  },
  "error": ""
}
```


## GetBech32Prefix

URL: [GET] /auth/bech32

### Response Example
```json
{
  "code": 200,
  "msg": {
    "bech32_prefix": "story"
  },
  "error": ""
}
```


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


## GetSupply

URL: [GET] /bank/supply

### Query Params
| Name | Type | Example | Require |
| --- | --- | --- | --- |
| pagination.key | string |  | ✔ |
| pagination.offset | string | 0 | ✔ |
| pagination.limit | string | 10 | ✔ |
| pagination.count_total | string | true | ✔ |
| pagination.reverse | string | true | ✔ |

### Response Example
```json
{
  "code": 200,
  "msg": {
    "supply": [
      {
        "denom": "stake",
        "amount": "79748037492401729"
      }
    ],
    "pagination": {
      "total": "1"
    }
  },
  "error": ""
}
```


## GetSupplyByDenom

URL: [GET] /bank/supply/by_denom

### Query Params
| Name | Type | Example | Require |
| --- | --- | --- | --- |
| denom | string | stake | ✔ |

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


## GetBalancesByAddress

URL: [GET] /bank/balances/{address}

### Path Params
| Name | Type | Example | Require |
| --- | --- | --- | --- |
| address | string | story1f5zuqhmwy39cv64g6laeeg264ydz06txlfqtg3 | ✔ |

### Query Params
| Name | Type | Example | Require |
| --- | --- | --- | --- |
| pagination.key | string |  | ✔ |
| pagination.offset | string | 0 | ✔ |
| pagination.limit | string | 10 | ✔ |
| pagination.count_total | string | true | ✔ |
| pagination.reverse | string | true | ✔ |

### Response Example
```json
{
  "code": 200,
  "msg": {
    "balances": [
      {
        "denom": "stake",
        "amount": "50301746"
      }
    ],
    "pagination": {
      "total": "1"
    }
  },
  "error": ""
}
```


## GetBalancesByAddressDenom

URL: [GET] /bank/balances/{address}/by_denom

### Path Params
| Name | Type | Example | Require |
| --- | --- | --- | --- |
| address | string | story1f5zuqhmwy39cv64g6laeeg264ydz06txlfqtg3 | ✔ |

### Query Params
| Name | Type | Example | Require |
| --- | --- | --- | --- |
| denom | string | stake |  |

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


## GetDenomOwners

URL: [GET] /bank/denom_owners/{denom}

### Path Params
| Name | Type | Example | Require |
| --- | --- | --- | --- |
| denom | string | stake | ✔ |

### Query Params
| Name | Type | Example | Require |
| --- | --- | --- | --- |
| pagination.key | string |  | ✔ |
| pagination.offset | string | 0 | ✔ |
| pagination.limit | string | 10 | ✔ |
| pagination.count_total | string | true | ✔ |
| pagination.reverse | string | true | ✔ |

### Response Example
```json
{
  "code": 200,
  "msg": {
    "denom_owners": [
      {
        "address": "story1le74erycve0n78lg6ylunmkp886dwxl26qgz80",
        "balance": {
          "denom": "stake",
          "amount": "15796442"
        }
      },
      {
        "address": "story1lxp43exras9rkk6kpqxq5phgv5ttd60el5vfx9",
        "balance": {
          "denom": "stake",
          "amount": "5000050186041"
        }
      },
      {
        "address": "story170yj2r7selm3dc0kkej8yfp0nwtgt456hgt5md",
        "balance": {
          "denom": "stake",
          "amount": "3771552"
        }
      },
      {
        "address": "story17zuu6pujj4jfzzg35qxcrj4ee0jgd9xdftv6z2",
        "balance": {
          "denom": "stake",
          "amount": "14192251"
        }
      },
      {
        "address": "story17z6fh2fqdqpa9am3a8m7c024r954j7hwxdeayf",
        "balance": {
          "denom": "stake",
          "amount": "1543497"
        }
      },
      {
        "address": "story1a4zcqjltpwl3zfk72kx7m2pycflgcxatn44dy6",
        "balance": {
          "denom": "stake",
          "amount": "131562519"
        }
      },
      {
        "address": "story1af4dyn4q27a460lxqnynp3zp8lgahsvlp8lm8v",
        "balance": {
          "denom": "stake",
          "amount": "154360480"
        }
      },
      {
        "address": "story1ay0e9kc5huhnyaqd6tj0eam6y7x94wuxu7kxks",
        "balance": {
          "denom": "stake",
          "amount": "155508"
        }
      },
      {
        "address": "story1u60eyhv44w578au4lv982wjkcuhhhyssynzpew",
        "balance": {
          "denom": "stake",
          "amount": "6253"
        }
      },
      {
        "address": "story1u3jrf5xtlkykm44u7py6xfdmktjfpnwxn9hz7c",
        "balance": {
          "denom": "stake",
          "amount": "1698454"
        }
      }
    ],
    "pagination": {
      "next_key": "FOGvHYtaPBrZdP8kj3AEJcmHvuTw",
      "total": "99"
    }
  },
  "error": ""
}
```


## GetDenomOwnersByQuery

URL: [GET] /bank/denom_owners_by_query

### Query Params
| Name | Type | Example | Require |
| --- | --- | --- | --- |
| pagination.key | string | stake |  |
| pagination.offset | string | 0 | ✔ |
| pagination.limit | string | 10 | ✔ |
| pagination.count_total | string | true | ✔ |
| pagination.reverse | string | true | ✔ |
| denom | string | stake |  |

### Response Example
```json
{
  "code": 200,
  "msg": {
    "denom_owners": [
      {
        "address": "story1le74erycve0n78lg6ylunmkp886dwxl26qgz80",
        "balance": {
          "denom": "stake",
          "amount": "15796442"
        }
      },
      {
        "address": "story1lxp43exras9rkk6kpqxq5phgv5ttd60el5vfx9",
        "balance": {
          "denom": "stake",
          "amount": "5000050186041"
        }
      },
      {
        "address": "story170yj2r7selm3dc0kkej8yfp0nwtgt456hgt5md",
        "balance": {
          "denom": "stake",
          "amount": "3771552"
        }
      },
      {
        "address": "story17zuu6pujj4jfzzg35qxcrj4ee0jgd9xdftv6z2",
        "balance": {
          "denom": "stake",
          "amount": "14192251"
        }
      },
      {
        "address": "story17z6fh2fqdqpa9am3a8m7c024r954j7hwxdeayf",
        "balance": {
          "denom": "stake",
          "amount": "1543497"
        }
      },
      {
        "address": "story1a4zcqjltpwl3zfk72kx7m2pycflgcxatn44dy6",
        "balance": {
          "denom": "stake",
          "amount": "131562519"
        }
      },
      {
        "address": "story1af4dyn4q27a460lxqnynp3zp8lgahsvlp8lm8v",
        "balance": {
          "denom": "stake",
          "amount": "154360480"
        }
      },
      {
        "address": "story1ay0e9kc5huhnyaqd6tj0eam6y7x94wuxu7kxks",
        "balance": {
          "denom": "stake",
          "amount": "155508"
        }
      },
      {
        "address": "story1u60eyhv44w578au4lv982wjkcuhhhyssynzpew",
        "balance": {
          "denom": "stake",
          "amount": "6253"
        }
      },
      {
        "address": "story1u3jrf5xtlkykm44u7py6xfdmktjfpnwxn9hz7c",
        "balance": {
          "denom": "stake",
          "amount": "1698454"
        }
      }
    ],
    "pagination": {
      "next_key": "FOGvHYtaPBrZdP8kj3AEJcmHvuTw",
      "total": "99"
    }
  },
  "error": ""
}
```


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


## GetDistributionValidatorsByDelegatorAddress

URL: [GET] /distribution/delegators/{delegator_address}/validators

### Path Params
| Name | Type | Example | Require |
| --- | --- | --- | --- |
| delegator_address | string | story1f5zuqhmwy39cv64g6laeeg264ydz06txlfqtg3 | ✔ |

### Response Example
```json
{
  "code": 200,
  "msg": {
    "validators": [
      "storyvaloper1rk277ak9enq422sevnf2yjjc6mhdv44wl8u30p",
      "storyvaloper1wk04lv2e2egfdt25e94xvwszny2ww8znancs5z",
      "storyvaloper137yg6cc4gpvhdx5wt0q3g5pcjeljuqvs4eu8y2",
      "storyvaloper1hcx6v7887r7r7498ge3kxqxrmkt3ekpk7jj5c2",
      "storyvaloper1mwgwgn6x5waumhg4tgnuhkezhamnaq4385mrgv",
      "storyvaloper17973uudyv484cgmy4sd4kjwdgrh86vfwhh0nlw"
    ]
  },
  "error": ""
}
```


## GetDelegatorWithdrawAddressByDelegatorAddress

URL: [GET] /distribution/delegators/{delegator_address}/withdraw_address

### Path Params
| Name | Type | Example | Require |
| --- | --- | --- | --- |
| delegator_address | string | story1f5zuqhmwy39cv64g6laeeg264ydz06txlfqtg3 | ✔ |

### Response Example
```json
{
  "code": 200,
  "msg": {
    "withdraw_address": "story1f5zuqhmwy39cv64g6laeeg264ydz06txlfqtg3"
  },
  "error": ""
}
```


## GetDelegatorRewardsByDelegatorAddress

URL: [GET] /distribution/delegators/{delegator_address}/rewards

### Path Params
| Name | Type | Example | Require |
| --- | --- | --- | --- |
| delegator_address | string | story1f5zuqhmwy39cv64g6laeeg264ydz06txlfqtg3 | ✔ |

### Response Example
```json
{
  "code": 200,
  "msg": {
    "rewards": [
      {
        "validator_address": "storyvaloper1rk277ak9enq422sevnf2yjjc6mhdv44wl8u30p",
        "reward": [
          {
            "denom": "stake",
            "amount": "519497.220124590080000000"
          }
        ]
      },
      {
        "validator_address": "storyvaloper1wk04lv2e2egfdt25e94xvwszny2ww8znancs5z",
        "reward": [
          {
            "denom": "stake",
            "amount": "895692.368287744000000000"
          }
        ]
      },
      {
        "validator_address": "storyvaloper137yg6cc4gpvhdx5wt0q3g5pcjeljuqvs4eu8y2",
        "reward": [
          {
            "denom": "stake",
            "amount": "16641260.184676218624000000"
          }
        ]
      },
      {
        "validator_address": "storyvaloper1hcx6v7887r7r7498ge3kxqxrmkt3ekpk7jj5c2",
        "reward": [
          {
            "denom": "stake",
            "amount": "12466838.774874030080000000"
          }
        ]
      },
      {
        "validator_address": "storyvaloper1mwgwgn6x5waumhg4tgnuhkezhamnaq4385mrgv",
        "reward": [
          {
            "denom": "stake",
            "amount": "534073.128124518400000000"
          }
        ]
      },
      {
        "validator_address": "storyvaloper17973uudyv484cgmy4sd4kjwdgrh86vfwhh0nlw",
        "reward": [
          {
            "denom": "stake",
            "amount": "16994236.393259084134400000"
          }
        ]
      }
    ],
    "total": [
      {
        "denom": "stake",
        "amount": "48051598.069346185318400000"
      }
    ]
  },
  "error": ""
}
```


## GetDelegatorRewardsByDelegatorAddressValidatorAddress

URL: [GET] /distribution/delegators/{delegator_address}/rewards/{validator_address}

### Path Params
| Name | Type | Example | Require |
| --- | --- | --- | --- |
| delegator_address | string | story1f5zuqhmwy39cv64g6laeeg264ydz06txlfqtg3 | ✔ |
| validator_address | string | storyvaloper17973uudyv484cgmy4sd4kjwdgrh86vfwhh0nlw | ✔ |

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


## GetDistributionValidatorByValidatorAddress

URL: [GET] /distribution/validators/{validator_address}

### Path Params
| Name | Type | Example | Require |
| --- | --- | --- | --- |
| validator_address | string | storyvaloper17973uudyv484cgmy4sd4kjwdgrh86vfwhh0nlw | ✔ |

### Response Example
```json
{
  "code": 200,
  "msg": {
    "operator_address": "story17973uudyv484cgmy4sd4kjwdgrh86vfwecmj59",
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


## GetValidatorCommissionByValidatorAddress

URL: [GET] /distribution/validators/{validator_address}/commission

### Path Params
| Name | Type | Example | Require |
| --- | --- | --- | --- |
| validator_address | string | storyvaloper17973uudyv484cgmy4sd4kjwdgrh86vfwhh0nlw | ✔ |

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
| Name | Type | Example | Require |
| --- | --- | --- | --- |
| validator_address | string | storyvaloper17973uudyv484cgmy4sd4kjwdgrh86vfwhh0nlw | ✔ |

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
| Name | Type | Example | Require |
| --- | --- | --- | --- |
| validator_address | string | storyvaloper17973uudyv484cgmy4sd4kjwdgrh86vfwhh0nlw | ✔ |

### Query Params
| Name | Type | Example | Require |
| --- | --- | --- | --- |
| pagination.key | string |  | ✔ |
| pagination.offset | string | 0 | ✔ |
| pagination.limit | string | 10 | ✔ |
| pagination.count_total | string | true | ✔ |
| pagination.reverse | string | true | ✔ |

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


## GetAllValidatorOutstandingRewards

URL: [GET] /distribution/all_validators/outstanding_rewards

### Query Params
| Name | Type | Example | Require |
| --- | --- | --- | --- |
| from | integer | 100 |  |
| to | integer | 200 |  |
