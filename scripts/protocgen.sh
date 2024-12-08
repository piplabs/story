#!/bin/bash
set -eo pipefail

# Ensure we are in the root of the repo, so  ${pwd}/go.mod must exit
if [ ! -f go.mod ]; then
  echo "Please run this script from the root of the repository"
  exit 1
fi

cd client/proto

echo "Generating proto files"
proto_dirs=$(find ./ -path -prune -o -name '*.proto' -print0 | xargs -0 -n1 dirname | sort | uniq)
for dir in $proto_dirs; do
  echo "DIR: ${dir}"
  # TODO: check proto files has condition: `if grep -q "option go_package.*story" "$file" &> /dev/null; then`

  if [[ $dir == *"/keeper" ]]; then
    echo "- orm protos"
    buf generate --template buf.gen.orm.yaml --path "${dir}"
  fi

  if [[ $dir == *"/module" ]]; then
    echo "- pulsar protos"
    buf generate --template buf.gen.pulsar.yaml --path "${dir}"
  fi

  if [[ $dir == *"/types" ]]; then
    echo "- gogo protos"
    buf generate --template buf.gen.gogo.yaml --path "${dir}"
  fi
done

echo "Move generated gogo proto files to the right places"
cp -r ./github.com/piplabs/story/client/* ../
rm -rf ./github.com

echo "Manually move orm proto files to the right places"
mv ./story/evmengine/v1/keeper/*.go ../x/evmengine/keeper/

cd ../..
