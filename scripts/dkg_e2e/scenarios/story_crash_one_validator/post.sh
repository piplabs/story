#!/usr/bin/env bash
# story_crash_one_validator/post.sh — Restart story on validator 2.
source "$(dirname "${BASH_SOURCE[0]}")/../_common.sh"

echo "[story_crash_one_validator] Post: restarting story on validator 2."

_ssh_cmd 2 "sudo systemctl start story 2>/dev/null || true"
sleep 30

echo "[story_crash_one_validator] Post: validator 2 story restarted."
