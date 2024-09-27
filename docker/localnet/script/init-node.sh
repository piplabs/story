#!/usr/bin/env bash

echo "node reset..."
rm -rf /home/ec2-user
cp -r /config /home/ec2-user
geth --state.scheme "hash" --gcmode archive init --datadir="/home/ec2-user/geth/data" /home/ec2-user/geth/config/genesis.json
openssl rand -hex 32 > /home/ec2-user/geth/data/geth/jwtsecret
cat <<EOF > /home/ec2-user/story/data/priv_validator_state.json
{
  "height": "0",
  "round": 0,
  "step": 0
}
EOF

sed -i "s/10.22.33.1/$(nslookup validator1 | grep Address | head -n 2 | tail -n 1 | awk -F ': ' '{print $2}')/g" /home/ec2-user/story/config/config.toml
sed -i "s/10.22.33.2/$(nslookup validator2 | grep Address | head -n 2 | tail -n 1 | awk -F ': ' '{print $2}')/g" /home/ec2-user/story/config/config.toml
sed -i "s/10.22.33.3/$(nslookup validator3 | grep Address | head -n 2 | tail -n 1 | awk -F ': ' '{print $2}')/g" /home/ec2-user/story/config/config.toml
sed -i "s/10.22.33.4/$(nslookup validator4 | grep Address | head -n 2 | tail -n 1 | awk -F ': ' '{print $2}')/g" /home/ec2-user/story/config/config.toml
sed -i "s/10.22.22.1/$(nslookup bootnode1 | grep Address | head -n 2 | tail -n 1 | awk -F ': ' '{print $2}')/g" /home/ec2-user/story/config/config.toml
sed -i "s/10.22.22.1/$(nslookup bootnode1 | grep Address | head -n 2 | tail -n 1 | awk -F ': ' '{print $2}')/g" /home/ec2-user/geth/config/geth.toml

echo "node run..."
geth --config /home/ec2-user/geth/config/geth.toml --nodekey /home/ec2-user/geth/config/nodekey --metrics --metrics.addr 0.0.0.0 &
story run --home=/home/ec2-user/story &

wait -n
exit $?
