#!/usr/bin/env bash

# This scripts generates all the protobufs using 'buf generate'.
# Cosmos however uses a mix of gogo, pulsar and orm plugins for code generation.
# So we manually call buf generate for each type of plugin.

function bufgen() {
    TYPE=$1 # Either orm,pulsar,proto
    DIR=$2 # Path to dir containing protos to generate
    OUTPUT=$3 # Path to output dir

    # Skip if ${DIR}/*.proto does not exist
    if ! test -n "$(find "${DIR}" -maxdepth 1 -name '*.proto')"; then
      return
    fi

    echo "  ${TYPE}: ${DIR}"

    buf generate \
      --template="../scripts/buf.gen.${TYPE}.yaml" \
      --output="${OUTPUT}" \
      --path="${DIR}"
}

# Ensure we are in the root of the repo, so  ${pwd}/go.mod must exit
if [ ! -f go.mod ]; then
  echo "Please run this script from the root of the repository"
  exit 1
fi

echo "Generating pulsar protos for cosmos module config"

bufgen pulsar story/evmengine/v1/module.proto ../client/x/evmengine/module
bufgen pulsar story/evmstaking/v1/module.proto ../client/x/evmstaking/module

cd proto

echo "Generating gogo protos for cosmos types"

bufgen gogo story/evmengine/v1/genesis.proto ../client/x/evmengine/types
bufgen gogo story/evmengine/v1/params.proto ../client/x/evmengine/types
bufgen gogo story/evmengine/v1/tx.proto ../client/x/evmengine/types

bufgen gogo story/evmstaking/v1/genesis.proto ../client/x/evmstaking/types
bufgen gogo story/evmstaking/v1/evmstaking.proto ../client/x/evmstaking/types
bufgen gogo story/evmstaking/v1/params.proto ../client/x/evmstaking/types
bufgen gogo story/evmstaking/v1/query.proto ../client/x/evmstaking/types
bufgen gogo story/evmstaking/v1/tx.proto ../client/x/evmstaking/types

echo "Generating orm protos"

bufgen orm proto/story/evmengine/v1/keeper/evmengine.proto ../client/x/evmengine/keeper

cd ..