#!/usr/bin/env bash

if [ -z "$NETWORK" ]; then
    NETWORK="public-testnet"
fi

DATA_DIR="/data"
INIT_DONE_FILE="$DATA_DIR/.init_done"
CONFIG_SOURCE="/init-config/$NETWORK"

if [ ! -d "$CONFIG_SOURCE" ]; then
    echo "network '$NETWORK' config not exists."
    exit 1
fi

if [ -f "$INIT_DONE_FILE" ]; then
    echo "skip init, run node..."
else
    echo "init config..."
    cp -r "$CONFIG_SOURCE/"* "/data/"

    story-init --config-dir="$DATA_DIR"
    geth --state.scheme "hash" --gcmode archive init --datadir="$DATA_DIR/geth/data" $DATA_DIR/geth/config/genesis.json
    openssl rand -hex 32 > $DATA_DIR/geth/data/geth/jwtsecret

    touch "$INIT_DONE_FILE"
fi

echo "node run..."
geth --config $DATA_DIR/geth/config/geth.toml --nodekey $DATA_DIR/geth/config/nodekey --metrics --metrics.addr 0.0.0.0 &
story run --home=$DATA_DIR/story &

wait -n
exit $?