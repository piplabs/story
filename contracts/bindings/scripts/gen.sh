#!/usr/bin/env bash
# Generate bindings for solidity contracts

DIR=${DIR:-./bindings}
PKG=${PKG:-bindings}

# generate bindings for the given contract
# works on contract name of fully qualified path to the contract
# params:
#  $1: contract name (ex. 'IPTokenStaking', 'src/protocol/IPTokenStaking.sol:IPTokenStaking')
gen_binding() {
  contract=$1

  # strip path prefix, if used
  # ex src/protocol/IPTokenStaking.sol:IPTokenStaking => IPTokenStaking
  name=$(echo ${contract} | cut -d ":" -f 2)

  # convert to lower case to respect golang package naming conventions
  name_lower=$(echo ${name} | tr '[:upper:]' '[:lower:]')

  temp=$(mktemp -d)
  forge inspect ${contract} abi --json > ${temp}/${name}.abi
  forge inspect ${contract} bytecode > ${temp}/${name}.bin

  abigen \
    --abi ${temp}/${name}.abi \
    --bin ${temp}/${name}.bin \
    --pkg ${PKG} \
    --type ${name} \
    --out ${DIR}/${name_lower}.go
}


for contract in $@; do
  gen_binding ${contract}
done
