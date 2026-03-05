#!/usr/bin/env bash

mockgen_cmd="mockgen"
$mockgen_cmd -package mock -destination testutil/mock/grpc_server.go github.com/cosmos/gogoproto/grpc Server
$mockgen_cmd -package mock -destination testutil/mock/logger.go cosmossdk.io/log Logger
$mockgen_cmd -source=../client/x/evmengine/types/expected_keepers.go -package testutil -destination ../client/x/evmengine/testutil/expected_keepers_mocks.go
$mockgen_cmd -source=../client/x/evmstaking/types/expected_keepers.go -package testutil -destination ../client/x/evmstaking/testutil/expected_keepers_mocks.go
$mockgen_cmd -source=../client/x/mint/types/expected_keepers.go -package testutil -destination ../client/x/mint/testutil/expected_keepers_mocks.go

# mock keepers for upgrade handlers
# Automatically generates mocks for any upgrade package that contains expected_keepers.go
for keepers_file in ../client/app/upgrades/*/expected_keepers.go; do
  [ -f "$keepers_file" ] || continue

  upgrade_dir=$(dirname "$keepers_file")
  testutil_dir="${upgrade_dir}/testutil"
  mkdir -p "$testutil_dir"

  echo "Generating mocks for ${keepers_file}..."
  $mockgen_cmd \
    -source="$keepers_file" \
    -package testutil \
    -destination "${testutil_dir}/expected_keepers_mocks.go"
done
