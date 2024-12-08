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
  proto_files=$(find "${dir}" -maxdepth 1 -name '*.proto')
  for file in $proto_files; do
    if grep -q "option go_package.*story" "$file" &> /dev/null; then
      if [[ $file == *"/module/"* ]]; then
        echo "  Generating for ${file} (pulsar)"
        buf generate --template buf.gen.pulsar.yaml "$file"
      elif [[ $file == *"/keeper/"* ]]; then
        echo "  Generating for ${file} (orm)"
        buf generate --template buf.gen.orm.yaml "$file"
      else
        echo "  Generating for ${file} (gogo)"
        buf generate --template buf.gen.gogo.yaml "$file"
      fi
    fi
  done
done

echo "Move generated gogo proto files to the right places"
cp -r ./github.com/piplabs/story/client/* ../
rm -rf ./github.com

cd ../..
