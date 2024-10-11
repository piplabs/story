#!/usr/bin/env bash

# Runs `pnpm lint-check` for every unique foundry project derived from the list
# of files provided as arguments by pre-commit.

source scripts/install_foundry.sh

source .pre-commit/foundry_utils.sh

for dir in $(foundryroots $@); do
  echo "Running 'lint-fix and lint-full' in ./$dir"
  (cd $dir && pnpm lint-fix && pnpm lint-full)
done
