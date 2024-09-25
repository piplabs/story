#!/usr/bin/env bash

if [ -z "$NETWORK" ]; then
    NETWORK="iliad"
fi

DATA_DIR="$HOME/.story"
COSMOVISOR_FILE_NAME="public-testnet"
INIT_DONE_FILE="$DATA_DIR/.init_done"

export DAEMON_NAME="story"
export DAEMON_HOME="/home/ec2-user/story"

case $NETWORK in
    iliad)
        echo "Use iliad as network ..."
        COSMOVISOR_FILE_NAME="public-testnet"
        ;;
    *)
        echo "Invalid Network. Use iliad."
        exit 1
        ;;
esac

if [ -f "$INIT_DONE_FILE" ]; then
    echo "skip init, run node..."
else
    echo "init network..."
    mkdir -p /home/ec2-user/story
    mkdir $DAEMON_HOME/data $DAEMON_HOME/backup

    # init cosmovisor data
    wget -q https://story-geth-binaries.s3.us-west-1.amazonaws.com/cosmovisor/$COSMOVISOR_FILE_NAME.tar.gz -O cosmovisor_data.tar.gz
    tar -zxf cosmovisor_data.tar.gz && rm -f cosmovisor_data.tar.gz
    mv cosmovisor $DAEMON_HOME/

    # init story data
    $DAEMON_HOME/cosmovisor/genesis/bin/story init --network iliad

    touch "$INIT_DONE_FILE"
fi

echo "node run..."
geth --iliad --syncmode full --http --http.addr 0.0.0.0 --http.vhosts '*' --ws --ws.addr 0.0.0.0 --ws.origins '*' --rpc.txfeecap 0 &
sleep 5s
cosmovisor run run --home=$DATA_DIR/story &

wait -n
exit $?
