#!/usr/bin/env bash
# story_crash_one_validator/pre.sh — Stop story on validator 2 to simulate crash.
# Chain should continue with 2/3 validators (sufficient for consensus).
source "$(dirname "${BASH_SOURCE[0]}")/../_common.sh"

echo "[story_crash_one_validator] Pre: stopping story on validator 2."

_ssh_cmd 2 "sudo systemctl stop story 2>/dev/null || true"
sleep 3

echo "[story_crash_one_validator] Pre: validator 2 story stopped. Chain continues with 2/3 consensus."
