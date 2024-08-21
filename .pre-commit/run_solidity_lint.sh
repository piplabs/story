#!/usr/bin/env bash

# Runs `pnpm lint-check` for every unique foundry project derived from the list
# of files provided as arguments by pre-commit.

source scripts/install_foundry.sh

# import foundryroots
source .pre-commit/foundry_utils.sh

for dir in $(foundryroots $@); do
  echo "Running 'lint-check' in ./$dir"
  (cd $dir && pnpm lint-check)
done
