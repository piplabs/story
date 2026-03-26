# Devnet test steps

Devnet은 1개의 bootnode와 3개의 validators로 구성되어있으며, validators가 돌아가는 machine들은 SGX support가 되는 machine이야.
config나 key들은 모두 재사용해야하기 때문에 이거는 건들지마. 위치는 ~/.story/story/에 있어 각 instance마다.
story-kernel에 코드 변경이 있는 경우 알다시피 code commitment(mrenclave)가 변경되기 때문에 변경된 값으로 genesis state를 만들어주고, 이를
업데이트해줘야해. 업데이트해야줘야하는 것들은 내가 아래 상세 스텝에서 설명해줄게.

## Steps

### 0. 기존 프로세스들 다 종료

모든 머신에서 실행:

```bash
# 모든 서비스 종료
sudo systemctl stop story node-geth story-kernel

# 혹시 남아있는 story-kernel 프로세스 강제 종료
sudo kill -9 $(sudo lsof -t -i :50051) 2>/dev/null || true

# 종료 확인
sudo systemctl status story node-geth story-kernel
```

### 1. origin repo에 배포

변경사항이 있는경우 story, story-kernel을 origin repo에 배포해야해. 사용하는 브랜치는 아래와 같고, 아래 브랜치에 한해서 너가 커밋할 수 있도록 권한을
줄게. 나머지는 안돼. 변경사항이 없다면 skip해도 좋아.

- story: origin/dkg/hans-temp-test
- story-kernel: origin/hans/temp-test

```bash
# story 배포 (로컬에서)
cd ~/repos/storyprotocol/story
git push origin HEAD:dkg/hans-temp-test

# story-kernel 배포 (변경사항 있을 경우)
cd ~/repos/storyprotocol/story-kernel
git push origin HEAD:hans/temp-test
```

### 2. ssh로 원격접속

접속정보는 다음과 같아. 아래 키와 접속정보를 사용해서 접속해.

```bash
# jpe-cdr-bootnode1
ssh -i ~/.ssh/dkg_new.pem -o StrictHostKeyChecking=no -o ConnectTimeout=10 ubuntu@23.102.71.16
# jpe-cdr-validator1
ssh -i ~/.ssh/dkg_new.pem -o StrictHostKeyChecking=no -o ConnectTimeout=10 ubuntu@20.46.165.193
# jpe-cdr-validator2
ssh -i ~/.ssh/dkg_new.pem -o StrictHostKeyChecking=no -o ConnectTimeout=10 ubuntu@20.48.25.208
# jpe-cdr-validator3
ssh -i ~/.ssh/dkg_new.pem -o StrictHostKeyChecking=no -o ConnectTimeout=10 ubuntu@40.115.139.113
```

### 3. story-kernel fetch 후 build, sign, 그리고 실행하기

이걸 먼저 하는 이유는 code commitment를 먼저 알기 위함이야. validator 3대에서 실행:

```bash
# 이전 state 제거
rm -rf ~/.story-kernel/keys/ && rm -rf ~/.story-kernel/dkg_state/

# 최신 코드 가져오기
cd ~/story-kernel
git fetch origin && git reset --hard origin/hans/temp-test

# 빌드
make build-with-cpp

# Gramine sign & mrenclave 확인 (이 값을 기록해둘 것, 3대 모두 같은 값이어야 함)
make all-gramine

# story-kernel 실행 및 로그 확인
sudo systemctl daemon-reload && sudo systemctl start story-kernel
journalctl -u story-kernel -f
```

### 4. 새로운 genesis state 생성

3에서 print된 mrenclave 값을 GenerateAlloc.s.sol의 SGX_CODE_COMMITMENT에 넣고, 스크립트를 실행하자. 로컬에서 실행:

```bash
cd ~/repos/storyprotocol/story

# contracts/script/GenerateAlloc.s.sol의 SGX_CODE_COMMITMENT를 mrenclave 값으로 업데이트 후:
CHAIN_ID=1511 make gengen
# 결과: contracts/local-alloc.json
```

### 5. 생성된 genesis state를 story-geth의 genesis.json에 update

~/genesis-geth.json에 있는 "alloc"의 value를 4에서 생성된 파일의 내용으로 업데이트 해줘. 로컬에서 실행:

```bash
# jq로 자동 교체
jq --slurpfile alloc contracts/local-alloc.json '.alloc = $alloc[0]' ~/genesis-geth.json > ~/genesis-geth-new.json && mv ~/genesis-geth-new.json ~/genesis-geth.json
```

### 6. story genesis.json의 evmengine params인 execution block hash update

5에서 업데이트된 execution layer의 genesis 파일을 사용해서 genesis block hash를 구해서 consensus layer의 evmengine params를 업데이트
해줘야해. 로컬에서 실행:

```bash
# story-geth로 genesis block hash 계산
cd ~/repos/storyprotocol/story-geth
rm -rf ~/.geth-tmmp && ./build/bin/geth --datadir ~/.geth-tmmp init ~/genesis-geth.json && ./build/bin/geth --datadir ~/.geth-tmmp console
# > eth.getBlock(0).hash
# "0x<block_hash_hex>"  <- 이 값을 복사
# > exit

# hex -> base64 변환 (0x 제거 후)
echo -n "<block_hash_hex_without_0x>" | xxd -r -p | base64

# ~/genesis-node.json의 evmengine params에서 execution_block_hash를 위 base64 값으로 교체
```

### 7. copy genesis files to each machine

scp를 사용해서 각 machine에 story와 story-geth genesis.json파일을 복사해줘. 로컬에서 실행:

```bash
# 모든 머신에 genesis 파일 복사 (for loop)
for host in 23.102.71.16 20.46.165.193 20.48.25.208 40.115.139.113; do
  scp -i ~/.ssh/dkg_new.pem ~/genesis-geth.json ubuntu@${host}:/home/ubuntu/config/genesis-geth.json
  scp -i ~/.ssh/dkg_new.pem ~/genesis-node.json ubuntu@${host}:/home/ubuntu/.story/story/config/genesis.json
done
```

### 8. origin branch에서 최신상태 가져와서 다시 build

모든 머신에서 실행:

```bash
cd ~/story
git fetch origin && git reset --hard origin/dkg/hans-temp-test && go build -o story ./client
```

### 9. 이전 테스트에서 생성됐던 data 파일들 삭제 후 network 재시작

모든 머신에서 실행:

```bash
# 이전 데이터 삭제
sudo rm -rf \
  /home/ubuntu/.story/geth/data/geth/{chaindata,blobpool,nodes} \
  /home/ubuntu/.story/story/data/* \
  /home/ubuntu/.story/story/config/write-file-atomic-* \
  /home/ubuntu/.story/story/config/addrbook.json

# priv_validator_state 초기화
echo '{"height": "0", "round": 0, "step": 0}' > /home/ubuntu/.story/story/data/priv_validator_state.json

# geth 초기화
geth --state.scheme=hash init --datadir=/home/ubuntu/.story/geth/data /home/ubuntu/config/genesis-geth.json

# 서비스 시작 및 로그 확인
sudo systemctl start story node-geth
journalctl -fu story
```

### 10. story-kernel light client 설정 후 시작 (validator 3대에서)

네트워크가 돌아간 뒤에 story-kernel의 light client config를 설정해야 한다. 체인이 블록을 생성하고 있어야
trusted height/hash를 얻을 수 있으므로 Step 9 이후에 진행한다.

```bash
# 1) trusted block 정보 확인 (아무 노드에서)
HEIGHT=$(curl -s localhost:26657/status | jq -r '.result.sync_info.latest_block_height')
TRUSTED_HEIGHT=$((HEIGHT - 5))
TRUSTED_HASH=$(curl -s "localhost:26657/block?height=$TRUSTED_HEIGHT" | jq -r '.result.block_id.hash')
echo "trusted_height = $TRUSTED_HEIGHT"
echo "trusted_hash = \"$TRUSTED_HASH\""

# 2) light_client 디렉토리 생성 (반드시 먼저 해야 함 — Gramine passthrough 이슈)
mkdir -p ~/.story-kernel/light_client

# 3) config.toml에 light_client 섹션 추가
# 각 validator마다 witness_addrs를 자기를 제외한 나머지 3개 노드로 설정
cat > ~/.story-kernel/config.toml << EOF
log-level = "info"

[grpc]
listen_addr = ":50051"

[light_client]
chain_id = "story-68930"
rpc_addr = "http://localhost:26657"
primary_addr = "http://localhost:26657"
witness_addrs = ["http://<other_node1>:26657", "http://<other_node2>:26657", "http://<other_node3>:26657"]
trusted_height = $TRUSTED_HEIGHT
trusted_hash = "$TRUSTED_HASH"
EOF

# 4) Gramine sign (코드 변경이 있었다면 mrenclave가 바뀔 수 있으므로 재실행)
cd ~/story-kernel && make all-gramine

# 5) story-kernel 시작
sudo systemctl daemon-reload && sudo systemctl start story-kernel
journalctl -u story-kernel -f
```

witness_addrs 예시 (validator별):
- validator1 (20.46.165.193): `["http://23.102.71.16:26657", "http://20.48.25.208:26657", "http://40.115.139.113:26657"]`
- validator2 (20.48.25.208): `["http://23.102.71.16:26657", "http://20.46.165.193:26657", "http://40.115.139.113:26657"]`
- validator3 (40.115.139.113): `["http://23.102.71.16:26657", "http://20.46.165.193:26657", "http://20.48.25.208:26657"]`

## 모니터링 명령어

```bash
# story 로그 실시간 확인
journalctl -fu story

# story-kernel 로그 실시간 확인
journalctl -u story-kernel -f

# geth 로그 실시간 확인
journalctl -fu node-geth

# 현재 블록 높이 확인 (story)
curl -s localhost:26657/status | jq '.result.sync_info.latest_block_height'

# 현재 블록 높이 확인 (geth)
curl -s -X POST -H "Content-Type: application/json" --data '{"jsonrpc":"2.0","method":"eth_blockNumber","params":[],"id":1}' http://localhost:8545 | jq -r '.result' | xargs printf "%d\n"

# validator 상태 확인
curl -s localhost:26657/validators | jq '.result.validators[] | {address: .address, voting_power: .voting_power}'

# peer 연결 상태 확인
curl -s localhost:26657/net_info | jq '.result.n_peers'

# 서비스 상태 확인
sudo systemctl status story node-geth story-kernel
```

## 트러블슈팅

```bash
# app hash mismatch 발생 시 — 1블록 롤백 (state wipe 대신)
sudo systemctl stop story
./story rollback
sudo systemctl start story

# story-kernel이 안 죽을 때
sudo kill -9 $(sudo lsof -t -i :50051) 2>/dev/null || true

# systemd restart loop 방지 — 반드시 stop 먼저 후 data wipe
sudo systemctl stop story

# conflicting vote 에러 시 — 1블록 롤백
sudo systemctl stop story
./story rollback
sudo systemctl start story

# story-kernel "light_client.db.lock: no such file or directory" 에러
# → light_client 디렉토리가 없음. 생성 후 재시작:
mkdir -p ~/.story-kernel/light_client
sudo systemctl restart story-kernel

# story-kernel "failed to acquire file lock: function not implemented" 에러
# → Gramine SGX에서 flock 미지원. story-kernel 코드에 flock fallback 패치 필요.
# → 패치 후 make build-with-cpp && make all-gramine 재실행

# story-kernel "chain id should not be empty" 에러
# → config.toml에 [light_client] 섹션이 누락됨. Step 10 참고

# story "no kernel clients available for registration" 에러
# → story-kernel이 실행 중이지 않음. story-kernel 먼저 시작해야 DKG 진행 가능
```
