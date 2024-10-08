#!/usr/bin/env bash

DATA_DIR="/home/ec2-user"
INIT_DONE_FILE="$DATA_DIR/.init_done"

cp /binary/* /usr/local/bin

if [ -f "$INIT_DONE_FILE" ]; then
  echo "skip init, run node..."
else
  echo "node $TYPE reset..."
  rm -rf $DATA_DIR
  cp -r /config $DATA_DIR
  geth --state.scheme "hash" --gcmode archive init --datadir="$DATA_DIR/geth/data" $DATA_DIR/geth/config/genesis.json
  openssl rand -hex 32 > $DATA_DIR/geth/data/geth/jwtsecret

  if [ "$TYPE" = "validator" ]; then
    sed -i "s/10.22.33.1/$(nslookup $(hostname) | grep Address | head -n 2 | tail -n 1 | awk -F ': ' '{print $2}')/g" $DATA_DIR/story/config/config.toml
  fi

  sed -i "s/10.22.22.1/$(nslookup bootnode1 | grep Address | head -n 2 | tail -n 1 | awk -F ': ' '{print $2}')/g" $DATA_DIR/story/config/config.toml
  sed -i "s/10.22.22.1/$(nslookup bootnode1 | grep Address | head -n 2 | tail -n 1 | awk -F ': ' '{print $2}')/g" $DATA_DIR/geth/config/geth.toml
  touch $INIT_DONE_FILE
fi

echo "node $TYPE run..."
geth --config $DATA_DIR/geth/config/geth.toml --nodekey $DATA_DIR/geth/config/nodekey --metrics --metrics.addr 0.0.0.0 &
story run --home=$DATA_DIR/story &

wait -n
exit $?
