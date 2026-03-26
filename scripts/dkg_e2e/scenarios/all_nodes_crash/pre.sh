#!/usr/bin/env bash
# all_nodes_crash/pre.sh — Stop ALL story + kernel on all validators.
# Simulates catastrophic failure. Geth stays running to preserve chain data.
source "$(dirname "${BASH_SOURCE[0]}")/../_common.sh"

TOTAL="${DKG_VALIDATOR_COUNT:-3}"
echo "[all_nodes_crash] Pre: stopping all story + kernel on ${TOTAL} validators."

for i in $(seq 1 "$TOTAL"); do
  _ssh_cmd "$i" "sudo systemctl stop story-kernel 2>/dev/null; sudo systemctl stop story 2>/dev/null; true" &
done
wait
sleep 3

echo "[all_nodes_crash] Pre: all nodes stopped. Chain halted."
