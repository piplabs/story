#!/usr/bin/env bash
set -eo pipefail

if [ ! -f go.mod ]; then
  echo "Please run this script from the root of the repository"
  exit 1
fi

cd client/dkg

echo "Generating dkg proto files"
protoc \
  --proto_path=proto \
  --proto_path=../proto \
  --proto_path=$(go list -m -f '{{.Dir}}' github.com/cosmos/gogoproto) \
  --go_out=pb --go_opt=paths=source_relative \
  --go-grpc_out=pb --go-grpc_opt=paths=source_relative \
  proto/v1/*.proto

cd ../..
