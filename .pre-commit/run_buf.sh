#!/usr/bin/env bash

if ! which buf 1>/dev/null; then
  echo "Installing buf"
  go generate scripts/tools.go
fi

EXPECT=$(go list -f "{{.Module.Version}}" github.com/bufbuild/buf/cmd/buf)
ACTUAL="v$(buf --version)"
if [[ "${EXPECT}" != "${ACTUAL}" ]]; then
  echo "Updating buf"
  go generate scripts/tools.go
fi

echo "buf version: $(buf --version)"
echo "protoc-gen-go version: $(protoc-gen-go --version)"

# Generate DKG service proto files
./scripts/dkg-protocgen.sh

# Generate Cosmos module proto files
./scripts/module-protocgen.sh
cd client/proto && buf lint && cd ../../
