#!/usr/bin/env bash
# all_nodes_crash/post.sh — Restart ALL story + kernel on all validators.
source "$(dirname "${BASH_SOURCE[0]}")/../_common.sh"

TOTAL="${DKG_VALIDATOR_COUNT:-3}"
echo "[all_nodes_crash] Post: restarting story on all ${TOTAL} validators."

# Start story first (consensus layer)
for i in $(seq 1 "$TOTAL"); do
  _ssh_cmd "$i" "sudo systemctl start story 2>/dev/null || true" &
done
wait
sleep 10

echo "[all_nodes_crash] Post: story restarted. Starting kernels..."

# Then start kernels
for i in $(seq 1 "$TOTAL"); do
  start_kernel "$i"
done

echo "[all_nodes_crash] Post: all nodes restarted."
