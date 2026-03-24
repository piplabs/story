#!/usr/bin/env bash
# Restore: cancel upgrade or let it complete naturally
source "$(dirname "${BASH_SOURCE[0]}")/../_common.sh"

echo "[upgrade_scheduled] Post: upgrade will be cancelled in Go test teardown (CancelDKGUpgrade)."
echo "  If the upgrade already activated, no cancel needed."
