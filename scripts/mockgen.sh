#!/usr/bin/env bash

mockgen_cmd="mockgen"
$mockgen_cmd -package mock -destination testutil/mock/grpc_server.go github.com/cosmos/gogoproto/grpc Server
$mockgen_cmd -package mock -destination testutil/mock/logger.go cosmossdk.io/log Logger
$mockgen_cmd -source=client/x/evmengine/types/expected_keepers.go -package testutil -destination client/x/evmengine/testutil/expected_keepers_mocks.go
$mockgen_cmd -source=client/x/evmstaking/types/expected_keepers.go -package testutil -destination client/x/evmstaking/testutil/expected_keepers_mocks.go
$mockgen_cmd -source=client/x/signal/types/expected_keepers.go -package testutil -destination client/x/signal/testutil/expected_keepers_mocks.go
