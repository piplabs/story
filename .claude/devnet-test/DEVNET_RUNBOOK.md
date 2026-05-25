# Devnet Testing Runbook

> Operational runbook for deploying and testing Story + story-kernel on the SGX devnet.
> Covers two workflows: **initial devnet setup** and **kernel upgrade (code commitment change)**.

## Branch Management

The `dkg/hans-temp-test` branch is the devnet deployment branch. It is always exactly
**one commit ahead** of `dkg/dev` — the single `feat: setup devnet` commit that contains
all devnet-specific configuration (GenerateAlloc, V200 height, DKG params, DCAP scripts).

When `dkg/dev` gets new commits, rebase `dkg/hans-temp-test` onto it:
```bash
git checkout dkg/hans-temp-test
git fetch origin
git rebase origin/dkg/dev   # reapplies "feat: setup devnet" on top
git push origin dkg/hans-temp-test --force
```

Any devnet configuration changes should be amended into the existing commit:
```bash
# make changes...
git add -A && git commit --amend --no-edit
git push origin dkg/hans-temp-test --force
```

---

## Infrastructure

| Role       | Host                | Public IP       | Internal IP | SGX | Services                        |
|------------|---------------------|-----------------|-------------|-----|---------------------------------|
| bootnode1  | jpe-hans-bootnode1  | 4.215.203.165   | 10.1.0.6    | No  | story, node-geth                |
| validator1 | jpe-hans-validator1 | 4.241.82.60     | 10.1.0.5    | Yes | story, node-geth, story-kernel  |
| validator2 | jpe-hans-validator2 | 74.176.184.25   | 10.1.1.5    | Yes | story, node-geth, story-kernel  |
| validator3 | jpe-hans-validator3 | 40.115.246.245  | 10.1.0.4    | Yes | story, node-geth, story-kernel  |

> **Internal IPs** are required for kernel light client `witness_addrs` (see Step 8.2).
> Find them with: `ssh <host> "hostname -I | awk '{print \$1}'"``

### SSH Access

Add the following to `~/.ssh/config` on your local machine (one-time setup):

```
Host boot
  HostName 4.215.203.165
  User ubuntu
  IdentityFile ~/.ssh/dkg_new.pem
  StrictHostKeyChecking no

Host val1
  HostName 4.241.82.60
  User ubuntu
  IdentityFile ~/.ssh/dkg_new.pem
  StrictHostKeyChecking no

Host val2
  HostName 74.176.184.25
  User ubuntu
  IdentityFile ~/.ssh/dkg_new.pem
  StrictHostKeyChecking no

Host val3
  HostName 40.115.246.245
  User ubuntu
  IdentityFile ~/.ssh/dkg_new.pem
  StrictHostKeyChecking no
```

Then connect with: `ssh val1`, `ssh val2`, `ssh val3`, `ssh boot`

For running commands across all machines:

```bash
for host in boot val1 val2 val3; do echo "=== $host ==="; ssh $host "<command>"; done
# Validators only:
for host in val1 val2 val3; do echo "=== $host ==="; ssh $host "<command>"; done
```

### Key Directories (on remote machines)

| Path                           | Description                            |
|--------------------------------|----------------------------------------|
| `~/.story/story/config/`      | CL config, genesis, validator keys     |
| `~/.story/story/data/`        | CL chain data                          |
| `~/.story/geth/data/`         | EL data                                |
| `~/.story-kernel/`            | Kernel state, config, keys             |
| `~/.story-kernel/light_client/` | Kernel light client data             |
| `~/config/genesis-geth.json`  | EL genesis file (source of truth)      |
| `~/story/`                    | story source code                      |
| `~/story-kernel/`             | story-kernel source code (old binary)  |

### Service Names

- `story` — consensus layer (CometBFT + app)
- `node-geth` — execution layer (NOT `geth`)
- `story-kernel` — SGX TEE kernel (systemd or manual)

### Critical Rules

**Identity & State**
- **NEVER delete** `priv_validator_key.json` or `node_key.json` — validator identities reused across resets.
- **NEVER wipe `~/.story/story/data/` after the chain starts** — destroys all on-chain state (DCAP, DKG). Requires full reset from Step 6.

**Branch & Binary**
- Always use `git checkout -f <branch>` (not `git reset --hard` alone) — reset doesn't switch branches.
- After branch switch, **rebuild + install binary**: `sudo systemctl stop story; make build; sudo cp build/story /usr/local/bin/story`. Binary cannot be replaced while story is running ("Text file busy").
- Verify correct binary after install: `story version | grep Commit` must match `git log --oneline -1`.

**Genesis**
- **Distribute genesis to ALL machines** before Step 6. Stale genesis on even one machine causes chain fork.
- After geth init, **verify genesis hash matches across all nodes**: `eth_getBlockByNumber("0x0")` must return identical hash.
- `execution_block_hash` goes in `app_state.evmengine.params.execution_block_hash` (NOT top-level).
- **EL chain_id** = `90930` (set in `genesis-geth.json` config and `story.toml` engine-chain-id). `1511` is the Foundry simulation chain ID used only for `forge script --chain-id` during alloc generation.
- **CL chain_id** = `story-68931` (must match `DKGTestChainID` in code).

**Alloc Generation**
- After `forge script GenerateAlloc.s.sol`, **manually add CREATE2 deployer** to `local-alloc.json` (see Step 5.2). The script's `vm.etch` doesn't persist to the alloc dump.
- `SGX_CODE_COMMITMENT` in GenerateAlloc.s.sol must match kernel MRENCLAVE.

**Process Management**
- **Step 1**: Use `sudo killall -9 loader gramine-sgx story geth` — SGX `loader` resists SIGTERM. Verify with `sudo lsof -i :50051`.
- **Step 7 geth init**: Run `sudo systemctl stop node-geth` **immediately before** `geth init` — systemd may auto-restart geth between Step 1 and Step 7, holding the datadir lock.
- **Step 6 data cleanup**: Wipe entire `~/.story/story/data/` directory (not individual files) to avoid store version mismatch ("initial version set to N, but found earlier version 1").
- **Kernel restart**: Never stop story to start kernel. Start kernel first, then `systemctl restart story`.

**DCAP**
- Use `DeployAllDevnet.s.sol` (NOT `deploy-dcap-devnet.sh` step-by-step) — deploys P256 + full DCAP stack atomically.
- **NEVER use `--skip-simulation`** with `forge script DeployAllDevnet`. It causes late transactions (V3QuoteVerifier, setQuoteVerifier, setAuthorized) to silently not broadcast. The V3QuoteVerifier will have no code and DCAP verification will fail with empty revert `0x`. Always run: `forge script ... --broadcast --legacy --gas-price 30000000000 -vv`.
- **daimo P256Verifier**: To use the daimo P256 verifier at `0xc2b78104907F722DABAc4C69f826a522B2754De4`, deploy it via CREATE2 deployer FIRST (see Step 7.5.0), then pass `P256_VERIFIER=0xc2b78104907F722DABAc4C69f826a522B2754De4` to `DeployAllDevnet`. Once deployed via CREATE2 on the first run, the P256Verifier persists in the genesis alloc and does not need re-deployment on subsequent resets (as long as the alloc is not regenerated).
- After DeployAllDevnet, **update `deployment-addresses.json`** with the new addresses. `upload-collateral.sh` reads from this file. Also update `automata-on-chain-pccs/deployment/90930.json` if using PCCS tools.
- **Verify V3QuoteVerifier has code** after deployment: `cast code <V3QV_ADDR> --rpc-url $RPC | wc -c` must be >100. If 3 (empty), the deployment was incomplete — re-run DeployAllDevnet WITHOUT `--skip-simulation`.
- **DAO authorization on DaoStorage**: `setCallerAuthorization` must be called for ALL 5 DAOs: PcsDao, PckDao, TcbEvalDao, EnclaveIdentityDaoVersioned, FmspcTcbDaoVersioned. DeployAllDevnet only calls `grantDao` (write access), but `setCallerAuthorization` (read access) is also required for collateral upload and quote verification.
- **SGX Hook address**: `0x93A1A0ce96A50BfBB858d34490Aee07BF18EA7f7` (deployed by `GenerateAlloc.s.sol` via Create3). Run `forge script GenerateAlloc.s.sol -vvv | grep SGXValidationHook` to confirm. The address `0xFC78D728...` in older configs is WRONG.
- **SGX Hook requires TWO calls** after each DCAP deploy: `setAutomataValidationAddr(address)` AND `setTcbEvaluationDataNumber(uint32)`. Missing `setTcbEvaluationDataNumber(18)` causes `TcbEvalExpiredOrNotFound` error during DKG registration.
- **Whitelist enclave type** must use the CORRECT SGX Hook proxy address (from GenerateAlloc output, NOT hardcoded).
- `cast --to-bytes32 1` may produce wrong padding. Use explicit: `0x0000000000000000000000000000000000000000000000000000000000000001`.
- All `cast send` requires `--legacy --gas-price 30000000000` (Story-geth doesn't support EIP-1559). **Do NOT store flags in a shell variable** (e.g., `F="--legacy --gas-price ..."`; use them inline).
- **tcbEvaluationDataNumber**: Must match what Intel PCS API returns (currently **18**, not 21). `deploy-dcap-devnet.sh` auto-detects this.
- DCAP + collateral can be deployed anytime — DKG rounds that miss the registration period simply skip and retry on the next round.

**Collateral Upload**
- `upload-collateral.sh` uploads certs in the WRONG ORDER: Platform CA before Root CA. Platform CA upload succeeds (Root CA not needed for CREATE), but subsequent Root CA shows `Duplicate_Collateral` if already in genesis. The real problem: **PCK CRL upload fails** because it verifies the CRL signature against Platform CA, which in turn needs Root CA to be present first.
- **Fix**: After `upload-collateral.sh`, always manually upload Platform CA cert and PCK CRL in order:
  1. `cast send $PCS_DAO 'upsertPcsCertificates(uint8,bytes)' 2 "0x$(openssl x509 -in .collateral-cache/issuer_cert_1.pem -outform DER | xxd -p | tr -d '\n')" ...`
  2. `cast send $PCS_DAO 'upsertPckCrl(uint8,bytes)' 2 "0x$(openssl crl -in .collateral-cache/pck_crl.pem -outform DER | xxd -p | tr -d '\n')" ...`
- **Verify PCK CRL is present**: `cast call $PCS_DAO 'getCertificateById(uint8)(bytes,bytes)' 2 --rpc-url $RPC` — the second return value (CRL) must NOT be `0x` (empty).

**Kernel Config Path**
- **Gramine manifest maps `--home /opt/story-kernel`** (NOT `~/.story-kernel`). Config file MUST be at `/opt/story-kernel/config.toml` on the host. Writing to `~/.story-kernel/config.toml` has no effect.
- Similarly, kernel state (light_client, keys, dkg_state) is at `/opt/story-kernel/` not `~/.story-kernel/`.

**Engine Chain ID**
- **`engine-chain-id` in `story.toml` must be `90930`** on ALL machines (validators AND bootnode). A common mistake is leaving it at `1511` (Foundry simulation chain ID). Verify after any reset: `grep engine-chain-id ~/.story/story/config/story.toml`.

**DKG Registration Period**
- Default `registration_period=50` (~2.5 min at 3s/block) is TOO SHORT to deploy DCAP + start kernels before Round 1 expires. Set `registration_period=500` (25 min) in the CL genesis (`app_state.dkg.params.registration_period`) to allow time for all setup.
- After modifying CL genesis, distribute to ALL machines before Step 6.

**Networking**
- **Kernel light client `witness_addrs` MUST use internal IPs** (10.1.x.x). SGX enclaves cannot route via public IPs.
- **Geth peers**: Boot node has 0 geth peers by default. Transactions sent to boot's RPC **never reach validators** and silently hang. **Always use a validator's RPC** (e.g., `http://4.241.82.60:8545`) for deployments, NOT the bootnode.
- For local `cast send` / `forge script`: set `DEVNET_RPC_URL` to a validator's **public IP:8545**.

---

# Part 1: Initial Devnet Setup

This section covers setting up a fresh devnet from scratch: building all components, generating genesis, and running the first DKG round.

## Step 1: Stop All Processes

Run on **all 4 machines**:

```bash
sudo systemctl stop story node-geth
# Validators only — stop kernel and kill ALL related processes:
sudo systemctl stop story-kernel 2>/dev/null
sleep 2
# Kill any processes holding kernel ports (SGX loader ignores SIGTERM)
sudo kill -9 $(sudo lsof -t -i :50051) 2>/dev/null
sudo kill -9 $(sudo lsof -t -i :50052) 2>/dev/null
# Kill ALL gramine/loader processes by exact name (-x, NOT -f)
sudo pkill -9 -x gramine-sgx 2>/dev/null
sudo pkill -9 -x loader 2>/dev/null
sleep 1
# Second pass — SGX cgroup cleanup can fail silently, leaving zombies
sudo kill -9 $(sudo lsof -t -i :50051) 2>/dev/null
sudo kill -9 $(sudo lsof -t -i :50052) 2>/dev/null
sudo pkill -9 -x gramine-sgx 2>/dev/null
sudo pkill -9 -x loader 2>/dev/null
```

**Verify (MANDATORY — do not proceed until clean):**

```bash
# Must show "inactive" for all
systemctl is-active story node-geth story-kernel 2>/dev/null || true
# Must show "kernel stopped" (no matching processes)
ps aux | grep -E 'gramine|story-kernel|loader' | grep -v grep || echo "kernel stopped"
# Must show NO listeners on kernel ports
sudo lsof -i :50051 -i :50052 2>/dev/null || echo "ports free"
```

> **If any process or port is still occupied:** run `sudo kill -9 <PID>` manually and re-verify.
> SGX enclaves sometimes resist cleanup due to cgroup kill failures ("Invalid argument").
> The two-pass kill above handles most cases, but always verify before proceeding.

## Step 2: Pull Latest Code on Remote Machines

> **IMPORTANT**: Use `git checkout -f`, NOT just `git reset --hard`.
> `git reset --hard origin/<branch>` moves HEAD but does NOT switch branches.
> The machine may silently stay on the wrong branch.

On **all 4 machines** (story):

```bash
cd ~/story
git fetch origin
git checkout -f <branch>           # e.g., dkg/hans-temp-test
git reset --hard origin/<branch>
git log --oneline -1               # VERIFY: must show expected commit
```

On **validators only** (story-kernel):

```bash
cd ~/story-kernel
git fetch origin
git checkout -f <branch>           # e.g., hans/audit-fix
git reset --hard origin/<branch>
git log --oneline -1               # VERIFY
```

## Step 3: Build story-kernel (Validators Only)

```bash
cd ~/story-kernel
make clean
make build-with-cpp
make all-gramine
```

After the build completes, note the **mrenclave** (code commitment) value from the output.

**Critical:** All 3 validators MUST produce the **same mrenclave** value. If they differ:
- Verify same commit hash on all machines (`git log --oneline -1`)
- Clean and rebuild (`make clean && make build-with-cpp && make all-gramine`)

### Enclave Memory

Gramine manifest requires `enclave_size >= 4G` for Go runtime. If you see `too many pages allocated` or
`failed to reserve page summary memory`, increase `enclave_size` in `story-kernel.manifest.template` and re-sign.
**Changing enclave_size changes the mrenclave.**

## Step 4: Build Story Consensus Client (All Machines)

On **all 4 machines**:

```bash
cd ~/story
make build

# Install binary — MUST stop story first (cp fails with "Text file busy" if running)
sudo systemctl stop story 2>/dev/null
sudo cp build/story /usr/local/bin/story

# VERIFY binary version matches the branch
story version | grep 'Git Commit'
# Must match: git log --oneline -1
```

## Step 5: Generate Genesis State (Local Machine)

This step runs on your **local dev machine**, not the remote servers.

### 5.1: Update GenerateAlloc Constants

Edit `contracts/script/GenerateAlloc.s.sol`:

```solidity
// Set mrenclave from Step 3
bytes32 constant SGX_CODE_COMMITMENT = 0x<mrenclave_hex>;

// CRITICAL for devnet: set to true so DKG, CDR, UBIPool, etc. are owned
// by the deployer key. This allows direct calls to whitelistEnclaveType,
// scheduleUpgrade, setMaxEncryptedDataSize, etc. without TimelockController.
bool private constant USE_DEPLOYER_AS_OWNER = true;
```

> **Why**: Without `USE_DEPLOYER_AS_OWNER = true`, all admin functions go through
> TimelockController. The timelock proposer key is set from `.env` (`ADMIN_ADDRESS`)
> and may not match the deployer key. This blocks HP-04 (kernel upgrade) and any
> runtime configuration changes on devnet.

### 5.2: Generate Alloc

```bash
cd story/contracts
forge script script/GenerateAlloc.s.sol --tc GenerateAlloc -vvv --chain-id 1511
# Output: alloc saved to ./local-alloc.json
```

**Post-processing: Add CREATE2 deployer and daimo P256Verifier to alloc.**
The `vm.etch` in GenerateAlloc doesn't persist the CREATE2 deployer to the alloc dump.
DCAP deploy scripts require it at `0x4e59b44847b379578588920cA78FbF26c0B4956C`.
The daimo P256Verifier at `0xc2b78104907F722DABAc4C69f826a522B2754De4` is required by
DCAP attestation for ECDSA P256 signature verification. Without it in genesis, cert
uploads and DKG registration will fail with empty revert `0x`.

```bash
cd story/contracts
python3 script/dcap-devnet/add-alloc-extras.py local-alloc.json
```

### 5.3: Update EL Genesis

Take the existing `genesis-geth.json` and replace ONLY the `"alloc"` field with contents from `local-alloc.json`. Keep all other fields (`config`, `difficulty`, `gasLimit`, etc.) unchanged.

You can do this with a script:

```bash
python3 -c "
import json
with open('genesis-geth.json') as f:
    genesis = json.load(f)
with open('contracts/local-alloc.json') as f:
    alloc = json.load(f)
genesis['alloc'] = alloc
with open('genesis-geth.json', 'w') as f:
    json.dump(genesis, f, indent=2)
print(f'Updated alloc with {len(alloc)} accounts')
"
```

> **First time?** Get genesis-geth.json from any existing machine: `scp val1:~/config/genesis-geth.json .`

### 5.4: Compute Execution Block Hash

**Important:** Use the same geth version as the remote machines (e.g., v1.2.1). Build it locally first.

```bash
cd story-geth
git checkout v1.2.1
make geth
rm -rf /tmp/geth-hash-tmp
./build/bin/geth --state.scheme=hash --datadir /tmp/geth-hash-tmp init path/to/genesis-geth.json
HASH=$(./build/bin/geth --datadir /tmp/geth-hash-tmp --exec "eth.getBlock(0).hash" console 2>/dev/null | tr -d '"')
echo $HASH
rm -rf /tmp/geth-hash-tmp
```

Take the returned hex hash (e.g., `0xabc123...`), base64-encode it:

```bash
echo -n "${HASH#0x}" | xxd -r -p | base64
```

### 5.5: Update CL Genesis

In the CL genesis file (`genesis-node.json`), set `app_state.evmengine.params.execution_block_hash` to the base64 value from above.

> **CRITICAL:** The `execution_block_hash` is located at `app_state.evmengine.params.execution_block_hash`, NOT as a top-level field. Putting it at the top level will silently fail.

Also update the CL genesis `chain_id` to match `DKGTestChainID` (currently `story-68931`).

### 5.6: Distribute Genesis Files

```bash
for host in boot val1 val2 val3; do
  scp genesis-geth.json $host:~/config/genesis-geth.json
  scp genesis-node.json $host:~/.story/story/config/genesis.json
done
```

> **CRITICAL**: The CL genesis source file is `genesis-node.json` (NOT `lib/netconf/local/genesis.json`
> which is a template with `{{LOCAL_ACCOUNT_ADDRESS}}` placeholders). Get the base file from any
> existing machine: `scp val2:~/config/genesis-node.json .`
>
> If even one machine has a stale genesis-geth.json, geth will initialize
> with a different genesis hash and the validator will fork from the network.
> After Step 7 (geth init), verify genesis hash match across all machines (see Step 7).

## Step 6: Clean Chain Data (All Machines)

Run on **all 4 machines**:

```bash
# Remove ALL chain data — wipe and recreate data dir from scratch.
# Individual file deletion can leave stale data that causes store version mismatches
# on restart (e.g., "initial version set to 2, but found earlier version 1").
sudo rm -rf ~/.story/geth/data/geth/{chaindata,blobpool,nodes,triecache}
sudo rm -rf ~/.story/story/data
mkdir -p ~/.story/story/data
sudo rm -f ~/.story/story/config/write-file-atomic-* ~/.story/story/config/addrbook.json

# Reset validator state (DO NOT DELETE the file — only reset contents)
echo '{"height": "0", "round": 0, "step": 0}' > ~/.story/story/data/priv_validator_state.json
```

On **validators only** — clean kernel state:

> **IMPORTANT**: Kernel state is at `/opt/story-kernel/` (Gramine mount), NOT `~/.story-kernel/`.
> Clean BOTH paths to be safe.

```bash
# Clean kernel state at BOTH possible paths
sudo rm -rf /opt/story-kernel/keys/ /opt/story-kernel/dkg_state/ /opt/story-kernel/light_client/ /opt/story-kernel/data/
rm -rf ~/.story-kernel/keys/ ~/.story-kernel/dkg_state/ ~/.story-kernel/light_client/ ~/.story-kernel/data/
```

On **validators only** — verify `story.toml` has correct kernel config:

```bash
grep -A5 '\[dkg\]' ~/.story/story/config/story.toml
grep kernel-endpoints ~/.story/story/config/story.toml
```

For a fresh setup with only one kernel, ensure:
```toml
[dkg]
enable = true
kernel-endpoints = ["127.0.0.1:50051"]
engine-rpc-endpoint = "http://127.0.0.1:8545"
enc-type = 1
```

Also, the chain ID of EL should be specified in `story.toml` if it is not specified.

```toml
engine-chain-id = 90930
```

> If `kernel-endpoints` contains extra entries (e.g., `127.0.0.1:50052` from a previous upgrade test), remove them unless you plan to run multiple kernels.

On **validators** — verify `priv_validator_key.json` matches the genesis validator set:

```bash
# Check your validator's CometBFT address:
python3 -c "import json; print(json.load(open('$HOME/.story/story/config/priv_validator_key.json'))['address'])"

# Compare with the on-chain validator set (from a running node, or check genesis):
# Your address MUST appear in the genesis validator set.
```

> **Warning:** If a previous test swapped validator keys (e.g., validator swap resharing), the key may be wrong.
> Look for backups: `ls ~/.story/story/config/priv_validator_key.json.bak*`

## Step 7: Initialize Geth and Start Services

> **IMPORTANT**: `node-geth` and `story` MUST be stopped before running `geth init`.
> If systemd has `Restart=always`, the service may auto-restart and hold a datadir lock,
> causing `geth init` to fail with "datadir already used by another process".

Run on **all 4 machines**:

```bash
# MUST stop services immediately before init — systemd Restart=always
# can re-launch geth between Step 1 and Step 7, holding the datadir lock.
sudo systemctl stop story node-geth 2>/dev/null
sleep 1

# Initialize geth with new genesis (MUST use --state.scheme=hash)
geth --state.scheme=hash init --datadir=$HOME/.story/geth/data $HOME/config/genesis-geth.json

# Start services
sudo systemctl daemon-reload
sudo systemctl start node-geth
sleep 2
sudo systemctl start story
```

**Verify blocks + genesis hash (MANDATORY):**

```bash
# Wait for blocks
sleep 15
journalctl -u story --no-pager -n 20 | grep "height="

# Verify genesis hash matches across ALL nodes (run from local):
for host in boot val1 val2 val3; do
  H=$(ssh $host "curl -s -X POST http://localhost:8545 \
    -H 'Content-Type: application/json' \
    -d '{\"jsonrpc\":\"2.0\",\"method\":\"eth_getBlockByNumber\",\"params\":[\"0x0\",false],\"id\":1}'" \
    | python3 -c "import json,sys; print(json.load(sys.stdin)['result']['hash'][:14])")
  echo "$host genesis=$H"
done
# ALL must show the SAME hash. If any differs, redistribute genesis and re-init that machine.
```

Wait until you see blocks being produced (height increasing).

## Step 7.5: Deploy DCAP and Upload Collateral

> **CRITICAL**: Must complete within the first DKG registration period (~200 blocks ≈ 6.5 min).
> After the registration period ends, DKG will skip the round and start a new one.
> DKG registration txs use the same deployer key and cause nonce collisions.

### Prerequisites (one-time per machine)

On **val1** (deployer machine):
```bash
cd ~
git clone https://github.com/automata-network/automata-on-chain-pccs.git
git clone https://github.com/automata-network/automata-dcap-attestation.git
cd ~/automata-on-chain-pccs && forge install
cd ~/automata-dcap-attestation && git submodule update --init --recursive
```

### 7.5.0: P256Verifier (included in genesis)

The daimo P256Verifier at `0xc2b78104907F722DABAc4C69f826a522B2754De4` is now included
in the genesis alloc via `add-alloc-extras.py` (Step 5.2). **No separate deployment needed.**

Verify after chain start:
```bash
cast code 0xc2b78104907F722DABAc4C69f826a522B2754De4 --rpc-url $RPC | wc -c
# Must show >100 (not 3 which means empty)
```

### 7.5.1: Deploy full DCAP stack with DeployAllDevnet

> **Use `DeployAllDevnet.s.sol`** instead of `deploy-dcap-devnet.sh`.
> This deploys P256 + all PCCS + DAOs + Router + Attestation + V3QuoteVerifier atomically
> and handles V3QV registration + Router authorization internally.
>
> If you deployed daimo P256Verifier in Step 7.5.0, set `P256_VERIFIER=0xc2b78104907F722DABAc4C69f826a522B2754De4`.
> Otherwise, DeployAllDevnet will deploy its own P256 verifier via CREATE.
>
> **CRITICAL: RPC must point to a VALIDATOR node, NOT the bootnode.**
> The bootnode's geth has 0 peers — transactions sent there never reach block producers
> and silently hang forever. Use `http://localhost:8545` when running from a validator,
> or the validator's public IP when running from your local machine.

```bash
KEY="0xac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80" #gitleaks:allow
# From val1: use localhost. From local: use val1's public IP.
RPC="http://localhost:8545"  # or http://4.241.82.60:8545 from local

cd ~/automata-dcap-attestation/evm

# ⚠️ CRITICAL: Clear old broadcast cache before each deploy.
# Stale cache from previous deployments causes nonce conflicts and
# "replacement transaction underpriced" errors that silently skip
# V3QuoteVerifier and grantDao transactions.
rm -rf broadcast/DeployAllDevnet.s.sol/$(cast chain-id --rpc-url $RPC)/ \
       cache_forge/DeployAllDevnet.s.sol/$(cast chain-id --rpc-url $RPC)/

OWNER=$(cast wallet address --private-key $KEY) \
P256_VERIFIER=0xc2b78104907F722DABAc4C69f826a522B2754De4 \
TCB_EVAL_NUMBER=18 \
forge script forge-script/DeployAllDevnet.s.sol:DeployAllDevnet \
    --rpc-url $RPC --private-key $KEY \
    --broadcast --legacy --gas-price 50000000000 -vv
# ⚠️ NEVER use --skip-simulation — it silently skips V3QuoteVerifier deployment
# ⚠️ Use 50 gwei gas-price (not 30) to avoid "replacement transaction underpriced"
```

**Post-deploy verification (MANDATORY):**
```bash
CHAIN_ID=$(cast chain-id --rpc-url $RPC)
python3 << 'PYEOF'
import json
txs = json.load(open(f"broadcast/DeployAllDevnet.s.sol/{CHAIN_ID}/run-latest.json"))["transactions"]
addrs = {}
for tx in txs:
    ca = tx.get("contractAddress")
    cn = tx.get("contractName","")
    h = tx.get("hash","")
    if ca and cn and cn not in addrs:
        addrs[cn] = (ca, "✓" if h else "⚠️ NO HASH")
for k,(v,s) in sorted(addrs.items()):
    print(f"  {s} {k}: {v}")
PYEOF
# ALL entries must show ✓. Any ⚠️ NO HASH means the tx was not broadcast.
# If V3QuoteVerifier shows ⚠️, re-run DeployAllDevnet WITHOUT --skip-simulation.
```

Expected output: `[15/15] V3QuoteVerifier registered and authorized` + address summary.

Extract deployed addresses from the broadcast file:
```bash
CHAIN_ID=$(cast chain-id --rpc-url $RPC)
python3 << 'PYEOF'
import json
txs = json.load(open(f"broadcast/DeployAllDevnet.s.sol/{CHAIN_ID}/run-latest.json"))["transactions"]
for tx in txs:
    ca = tx.get("contractAddress")
    cn = tx.get("contractName", "")
    if ca: print(f"  {cn}: {ca}")
PYEOF
```

### 7.5.2: Authorize DAOs on DaoStorage

DeployAllDevnet calls `grantDao()` (write access) internally, but these calls may silently fail
due to nonce conflicts (especially on re-deployments). Always run `grantDao` manually.
`setCallerAuthorization()` (read access) is also required. Authorize ALL 5 DAOs:

```bash
# Parse addresses from broadcast JSON
CHAIN_ID=$(cast chain-id --rpc-url $RPC)
eval $(python3 << 'PYEOF'
import json
txs = json.load(open(f"broadcast/DeployAllDevnet.s.sol/{CHAIN_ID}/run-latest.json"))["transactions"]
addrs = {}
for tx in txs:
    ca = tx.get("contractAddress")
    cn = tx.get("contractName","")
    if ca and cn and cn not in addrs: addrs[cn] = ca
print(f'DAO_STORAGE={addrs["AutomataDaoStorage"]}')
print(f'PCS_DAO={addrs["AutomataPcsDao"]}')
print(f'PCK_DAO={addrs["AutomataPckDao"]}')
print(f'TCB_EVAL_DAO={addrs["AutomataTcbEvalDao"]}')
print(f'ENCLAVE_DAO={addrs["AutomataEnclaveIdentityDaoVersioned"]}')
print(f'FMSPC_DAO={addrs["AutomataFmspcTcbDaoVersioned"]}')
print(f'ATTESTATION={addrs["AutomataDcapAttestationFee"]}')
PYEOF
)

# Grant DAO write access (grantDao) + read access (setCallerAuthorization)
for DAO in $PCS_DAO $PCK_DAO $TCB_EVAL_DAO $ENCLAVE_DAO $FMSPC_DAO; do
  cast send $DAO_STORAGE 'grantDao(address)' $DAO \
    --rpc-url $RPC --private-key $KEY --legacy --gas-price 50000000000
  cast send $DAO_STORAGE 'setCallerAuthorization(address,bool)' $DAO true \
    --rpc-url $RPC --private-key $KEY --legacy --gas-price 50000000000
done

# Verify: ALL must return true
for DAO in $PCS_DAO $PCK_DAO $TCB_EVAL_DAO $ENCLAVE_DAO $FMSPC_DAO; do
  echo "$DAO: $(cast call $DAO_STORAGE 'isAuthorizedCaller(address)(bool)' $DAO --rpc-url $RPC)"
done
```

### 7.5.3: Connect SGX Hook + Whitelist Enclave

> **CRITICAL**: The SGX Hook address is NOT `0xFC78D728...`. It is deployed by `GenerateAlloc.s.sol` via Create3.
> Find the correct address by running: `cd story/contracts && forge script script/GenerateAlloc.s.sol --tc GenerateAlloc -vvv --chain-id 1511 2>&1 | grep "SGXValidationHook proxy"`
> Current address: `0x93A1A0ce96A50BfBB858d34490Aee07BF18EA7f7`

```bash
MRENCLAVE="<mrenclave_from_step_3>"
SGX_HOOK="0x93A1A0ce96A50BfBB858d34490Aee07BF18EA7f7"  # From GenerateAlloc output (NOT 0xFC78D728...)
DKG="0xCcCcCC0000000000000000000000000000000004"
# $ATTESTATION from Step 7.5.2 (parsed from broadcast JSON)

# ⚠️ BOTH calls are required. Missing setTcbEvaluationDataNumber causes TcbEvalExpiredOrNotFound.
cast send $SGX_HOOK 'setAutomataValidationAddr(address)' $ATTESTATION \
    --rpc-url $RPC --private-key $KEY --legacy --gas-price 30000000000
cast send $SGX_HOOK 'setTcbEvaluationDataNumber(uint32)' 18 \
    --rpc-url $RPC --private-key $KEY --legacy --gas-price 30000000000

# Whitelist enclave — explicit bytes32, do NOT use cast --to-bytes32
ENCLAVE_TYPE="0x0000000000000000000000000000000000000000000000000000000000000001"
cast send $DKG 'whitelistEnclaveType(bytes32,(bytes32,address),bool)' \
    $ENCLAVE_TYPE "(0x$MRENCLAVE,$SGX_HOOK)" true \
    --rpc-url $RPC --private-key $KEY --legacy --gas-price 30000000000

# Verify (ALL must match expected values)
echo "automataValidationAddr: $(cast call $SGX_HOOK 'automataValidationAddr()(address)' --rpc-url $RPC)"
echo "tcbEvaluationDataNumber: $(cast call $SGX_HOOK 'tcbEvaluationDataNumber()(uint32)' --rpc-url $RPC)"
echo "enclaveWhitelisted: $(cast call $DKG 'isEnclaveTypeWhitelisted(bytes32)(bool)' $ENCLAVE_TYPE --rpc-url $RPC)"
# automataValidationAddr MUST equal $ATTESTATION
# tcbEvaluationDataNumber MUST be 18
# enclaveWhitelisted MUST be true
```

### 7.5.4: Upload Intel Collateral

```bash
# Grant ATTESTER_ROLE on versioned DAOs
# $ENCLAVE_DAO and $FMSPC_DAO from Step 7.5.2
OWNER=$(cast wallet address --private-key $KEY)
cast send $ENCLAVE_DAO 'grantRoles(address,uint256)' $OWNER 1 \
    --rpc-url $RPC --private-key $KEY --legacy --gas-price 30000000000
cast send $FMSPC_DAO 'grantRoles(address,uint256)' $OWNER 1 \
    --rpc-url $RPC --private-key $KEY --legacy --gas-price 30000000000

# Update deployment-addresses.json BEFORE running upload-collateral.sh
# (use addresses from Step 7.5.1 broadcast output)

# Upload collateral
cd ~/story/contracts/script/dcap-devnet
export FMSPC="00606A000000"  # Azure DCsv3 Ice Lake-SP
export DEVNET_RPC_URL="$RPC" DEPLOYER_PRIVATE_KEY="$KEY"
export PCCS_REPO=~/automata-on-chain-pccs DCAP_REPO=~/automata-dcap-attestation TCB_EVAL_NUMBER=18
./upload-collateral.sh
```

**Post-upload fix (MANDATORY):** `upload-collateral.sh` uploads Platform CA before Root CA,
causing PCK CRL upload to fail. After the script completes, manually fix:

```bash
cd ~/story/contracts/script/dcap-devnet
# $PCS_DAO from Step 7.5.2

# 1. Re-upload Platform CA cert (now Root CA is present for signature verification)
PLATFORM_CA_HEX=$(openssl x509 -in .collateral-cache/issuer_cert_1.pem -outform DER | xxd -p | tr -d '\n')
cast send $PCS_DAO 'upsertPcsCertificates(uint8,bytes)' 2 "0x${PLATFORM_CA_HEX}" \
    --rpc-url $RPC --private-key $KEY --gas-limit 5000000 --legacy --gas-price 30000000000

# 2. Upload PCK CRL (requires Platform CA to be present)
PCK_CRL_HEX=$(openssl crl -in .collateral-cache/pck_crl.pem -outform DER | xxd -p | tr -d '\n')
cast send $PCS_DAO 'upsertPckCrl(uint8,bytes)' 2 "0x${PCK_CRL_HEX}" \
    --rpc-url $RPC --private-key $KEY --gas-limit 5000000 --legacy --gas-price 30000000000
```

### 7.5.5: Verify

```bash
# 1. Check PCK CRL is present (CRITICAL — empty CRL causes attestation failure)
cast call $PCS_DAO 'getCertificateById(uint8)(bytes,bytes)' 2 --rpc-url $RPC 2>&1 | python3 -c "
import sys; lines=sys.stdin.readlines()
cert_len = len(lines[0].strip()) if len(lines) > 0 else 0
crl_len = len(lines[1].strip()) if len(lines) > 1 else 0
print(f'Platform CA cert: {\"✓\" if cert_len > 10 else \"MISSING\"} ({cert_len} chars)')
print(f'PCK CRL: {\"✓\" if crl_len > 10 else \"MISSING — MUST FIX\"} ({crl_len} chars)')
"

# 2. Check geth logs for DCAP errors
journalctl -u node-geth --no-pager | grep 'eth_estimateGas.*err' | tail -3
# SUCCESS: no output or estimateGas at DEBUG level (no err)
# FAILURE: errdata 0xa78bf21a → TcbEvalExpiredOrNotFound (missing setTcbEvaluationDataNumber or TCB Info not uploaded)
# FAILURE: errdata 0x482b7129 → CrlExpiredOrNotFound (PCK CRL missing)
# FAILURE: errdata "0x" (empty) → V3QuoteVerifier has no code (--skip-simulation issue), or DAO not authorized, or SGX Hook has no code
```

## Step 8: Configure and Start story-kernel (Validators Only)

story-kernel's light client needs a trusted block from the running chain. Wait for the chain to produce ~5 blocks first.

### 8.1: Get Trusted Block Info

Run from your **local machine** (avoids python quoting issues on remote):

```bash
# Wait for blocks first, then query a recent block
ssh val1 "curl -s 'http://localhost:26657/block?height=5'" | python3 -c "
import json, sys
r = json.load(sys.stdin)['result']
print(f'height: {r[\"block\"][\"header\"][\"height\"]}')
print(f'hash: {r[\"block_id\"][\"hash\"]}')
"
```

Note the `height` and `hash` — you'll need them for the kernel config.

### 8.2: Write Kernel Config

> **CRITICAL**: The Gramine manifest maps `--home /opt/story-kernel`. Config MUST be written to
> `/opt/story-kernel/config.toml` on the host (use `sudo tee`), NOT `~/.story-kernel/config.toml`.
> Writing to `~/.story-kernel/` has no effect — the enclave reads from `/opt/story-kernel/`.

Create/overwrite `/opt/story-kernel/config.toml` on each validator (use `sudo tee`):

```toml
log-level = "info"

[grpc]
listen_addr = ":50051"

[light_client]
chain_id = "story-68931"
rpc_addr = "http://localhost:26657"
primary_addr = "http://localhost:26657"
witness_addrs = ["http://<other_val1_ip>:26657", "http://<other_val2_ip>:26657"]
trusted_height = <height_from_above>
trusted_hash = "<hash_from_above>"
```

> **CRITICAL:** `witness_addrs` MUST use **internal IPs** (10.x.x.x), NOT public IPs.
> SGX enclaves on Azure VMs cannot route CometBFT RPC through public IPs.
> Point to OTHER validators (not self, not bootnode).
>
> Examples (using internal IPs from Infrastructure table):
> - val1 (10.1.0.5): `witness_addrs = ["http://10.1.1.5:26657", "http://10.1.0.4:26657"]`
> - val2 (10.1.1.5): `witness_addrs = ["http://10.1.0.5:26657", "http://10.1.0.4:26657"]`
> - val3 (10.1.0.4): `witness_addrs = ["http://10.1.0.5:26657", "http://10.1.1.5:26657"]`

### 8.3: Start Kernel

```bash
# If using systemd:
sudo systemctl start story-kernel

# If running manually (for debugging):
cd ~/story-kernel
nohup gramine-sgx story-kernel start --home ~/.story-kernel > /tmp/kernel.log 2>&1 &
```

> **Note:** Gramine SGX enclave loading takes 1-3 minutes. Do not panic if it appears stuck at "Parsing TOML manifest file".

### 8.4: Verify Kernel Running

Wait 1-3 minutes for enclave loading, then check the log:

```bash
# Manual mode:
tail -5 /tmp/kernel.log
# Should see: "gRPC server listening on :50051" and "VerifyHeader" entries

# Systemd mode:
journalctl -u story-kernel --no-pager -n 20
```

### 8.5: Restart Story to Connect to Kernel

After kernel starts, **restart story** so it picks up the kernel connection:

> **WARNING**: Use `systemctl restart story` ONLY — do NOT wipe `~/.story/story/data/`.
> Wiping data after chain start deletes all on-chain state (DCAP collateral, DKG rounds)
> and causes store version mismatch errors on restart.

```bash
sudo systemctl restart story
```

Verify kernel connection in story logs:

```bash
journalctl -u story --no-pager -n 100 | grep "Connected to kernel"
# Expected: Connected to kernel endpoint  endpoint=127.0.0.1:50051 code_commitment=<mrenclave>
```

## Step 9: Verify DKG Round Completion

The first DKG round starts automatically. Monitor progress:

```bash
# Watch for DKG lifecycle events
journalctl -u story -f | grep -E "DKG|Emitted"
```

Expected sequence:
1. `Initiated new DKG round  round=N` — round started
2. `Handling DKG registration` → `DKG initialization complete` — registered
3. `Emitted BeginDKGDealing event` — dealing phase started
4. `GenerateDeals call to kernel client` — deals generated
5. `Emitted BeginDKGFinalization event` — finalization phase started
6. `DKG successfully finalized` — each validator finalized
7. `Emitted DKGFinalized event  round=N` — round complete!
8. `DKG process completed successfully  round=N` — all done

If a round fails, check:
- `journalctl -u story --no-pager | grep -i "error\|failed\|mismatch"` for errors
- `journalctl -u story-kernel --no-pager | grep -i "error"` for kernel errors

---

# Part 2: Kernel Upgrade (Code Commitment Change)

This section covers testing the DKG upgrade resharing flow: deploying a new kernel binary with a different code
commitment (mrenclave), whitelisting it on-chain, scheduling an upgrade, and verifying that the DKG reshares keys to the
new binary.

**Prerequisites:** A working devnet with at least one completed DKG round (Part 1).

## Overview

The upgrade flow:
1. Build a new kernel binary → produces a **new code commitment**
2. Run the new binary alongside the old one (different port)
3. Update story config to connect to both kernels
4. Whitelist the new code commitment on-chain
5. Schedule the upgrade at a future block height
6. DKG automatically triggers upgrade resharing at the activation height
7. After resharing completes, the old kernel can be decommissioned

## Step 1: Build New Kernel Binary

On **each validator**, build the new kernel in a separate directory:

```bash
# Clone or copy to a new directory
cp -r ~/story-kernel ~/story-kernel-new
cd ~/story-kernel-new

# Checkout the new branch/version
git fetch origin
git checkout <new-branch>
git reset --hard origin/<new-branch>

# Build
make clean
make build-with-cpp
make all-gramine
```

Note the new mrenclave value. It MUST differ from the old kernel's mrenclave.

**Verify both mrenclave values are different:**

```bash
echo "Old kernel mrenclave:"
cat ~/story-kernel/story-kernel.manifest.sgx.d/mrenclave.txt 2>/dev/null || echo "check build output"
echo "New kernel mrenclave:"
cat ~/story-kernel-new/story-kernel.manifest.sgx.d/mrenclave.txt 2>/dev/null || echo "check build output"
```

## Step 2: Start New Kernel on Port 50052

On **each validator**, start the new kernel alongside the old one:

### 2.1: Create Config for New Kernel

```bash
mkdir -p ~/.story-kernel-new

# Copy and modify config — use port 50052 instead of 50051
cat > ~/.story-kernel-new/config.toml << 'EOF'
log-level = "info"

[grpc]
listen_addr = ":50052"

[light_client]
chain_id = "<chain_id>"
rpc_addr = "http://localhost:26657"
primary_addr = "http://localhost:26657"
witness_addrs = ["http://<other_val1_ip>:26657", "http://<other_val2_ip>:26657"]
trusted_height = <recent_height>
trusted_hash = "<recent_hash>"
EOF
```

> Use a recent block height/hash from the running chain (see Part 1, Step 8.1).

### 2.2: Start New Kernel

```bash
cd ~/story-kernel-new
nohup gramine-sgx story-kernel start --home ~/.story-kernel-new > /tmp/kernel-new.log 2>&1 &
```

### 2.3: Verify Both Kernels Running

```bash
# Old kernel on 50051
curl -s http://localhost:50051/health 2>/dev/null && echo "old kernel OK" || echo "old kernel FAIL"

# New kernel on 50052
curl -s http://localhost:50052/health 2>/dev/null && echo "new kernel OK" || echo "new kernel FAIL"

# Or check processes
pgrep -f "story-kernel" -a
```

## Step 3: Update Story Config to Connect to Both Kernels

On **each validator**, edit `~/.story/story/config/story.toml`:

```toml
[dkg]
enable = true
kernel-endpoints = ["127.0.0.1:50051", "127.0.0.1:50052"]
```

Then restart story:

```bash
sudo systemctl restart story
```

Verify both kernel connections:

```bash
journalctl -u story --no-pager -n 100 | grep "Connected to kernel"
# Expected: two lines, one for each endpoint with different code_commitment values
```

Example output:
```
Connected to kernel endpoint  endpoint=127.0.0.1:50051 code_commitment=<OLD_CC>
Connected to kernel endpoint  endpoint=127.0.0.1:50052 code_commitment=<NEW_CC>
```

## Step 4: Whitelist New Code Commitment On-Chain

**This MUST be done BEFORE scheduling the upgrade.**

### 4.1: Get Validation Hook Address

Query the existing enclave type 1 for its validation hook:

```bash
# On any validator:
PATH=$PATH:$HOME/.foundry/bin

cast call 0xCcCcCC0000000000000000000000000000000004 \
  "enclaveTypeData(bytes32)" \
  $(cast --to-bytes32 1) \
  --rpc-url http://localhost:8545
```

The response contains two 32-byte slots: `[codeCommitment, validationHook]`. The validation hook address is in the second slot (last 20 bytes).

### 4.2: Whitelist New Enclave Type

```bash
NEW_CC=<new_kernel_mrenclave_hex_no_0x_prefix>
HOOK_ADDR=<validation_hook_from_4.1>
DKG_OWNER_KEY=0xac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80 #gitleaks:allow

PATH=$PATH:$HOME/.foundry/bin
cast send 0xCcCcCC0000000000000000000000000000000004 \
  "whitelistEnclaveType(bytes32,(bytes32,address),bool)" \
  $(cast --to-bytes32 2) \
  "(0x${NEW_CC},${HOOK_ADDR})" \
  true \
  --private-key $DKG_OWNER_KEY \
  --rpc-url http://localhost:8545 \
  --gas-price 17000000000 \
  --legacy
```

### 4.3: Verify Whitelisting

```bash
cast call 0xCcCcCC0000000000000000000000000000000004 \
  "isEnclaveTypeWhitelisted(bytes32)" \
  $(cast --to-bytes32 2) \
  --rpc-url http://localhost:8545
```

Expected: `0x...01` (true). If `0x...00`, the whitelist tx failed — check the tx receipt.

## Step 5: Schedule the Upgrade

### 5.1: Check Current Block Height

```bash
journalctl -u story --no-pager -n 5 | grep "height="
```

### 5.2: Schedule Upgrade

Choose an activation height at least **50-100 blocks** in the future. The upgrade resharing round starts when the chain reaches this height.

```bash
ACTIVATION_HEIGHT=<current_height + 100>

PATH=$PATH:$HOME/.foundry/bin
cast send 0xCcCcCC0000000000000000000000000000000004 \
  "scheduleUpgrade(uint256,string)" \
  $ACTIVATION_HEIGHT \
  "v2.0.0" \
  --private-key $DKG_OWNER_KEY \
  --rpc-url http://localhost:8545 \
  --gas-price 17000000000 \
  --legacy
```

### 5.3: Verify Upgrade Scheduled

```bash
journalctl -u story --no-pager | grep -i "UpgradeScheduled\|upgrade.*scheduled"
```

## Step 6: Monitor Upgrade Resharing

### 6.1: Wait for Activation Height

```bash
journalctl -u story -f | grep -E "Kernel upgrade activated|initiating upgrade resharing"
```

### 6.2: Monitor Registration

```bash
journalctl -u story -f | grep -E "Handling DKG registration|Resolved enclave type|Register succeeded"
```

Key things to verify:
- `Resolved enclave type for registration  enclave_type=...0002` — must be the NEW type
- `Register succeeded` — registration tx mined

### 6.3: Monitor Dealing and Finalization

```bash
journalctl -u story -f | grep -E "deals|responses|Finalize|DKGFinalized|DKG process completed"
```

### 6.4: Verify Success

```bash
journalctl -u story --no-pager | grep -E "DKGFinalized|DKG process completed" | tail -5
```

Expected:
```
Emitted DKGFinalized event  round=N
DKG process completed successfully  round=N
```

### 6.5: Verify GlobalPubKey Preservation

Compare the GlobalPubKey from the upgrade round with the previous round:

```bash
journalctl -u story --no-pager | grep "global_pub_key=" | tail -5
```

The GlobalPubKey MUST be **identical** before and after resharing.

## Step 7: Post-Upgrade Cleanup (Optional)

After confirming the upgrade resharing succeeded:

1. The old kernel (port 50051) can be stopped
2. Update story config to remove the old endpoint
3. Optionally rename `story-kernel-new` to `story-kernel`

---

## Quick Reference: What Changes vs. What Stays When Resetting

| Item | Reset? | Notes |
|------|--------|-------|
| `priv_validator_key.json` | **NO** | Validator identity, reused across resets |
| `node_key.json` | **NO** | P2P identity, reused across resets |
| `priv_validator_state.json` | **RESET contents** | Reset to height 0, do NOT delete |
| `genesis.json` (CL) | **REPLACE** | New execution_block_hash |
| `genesis-geth.json` (EL) | **REPLACE** | New alloc with updated code commitment |
| Chain data (CL + EL) | **DELETE** | Fresh start from genesis |
| Kernel state (`~/.story-kernel/`) | **DELETE** (keys, dkg_state, light_client, data) | Old key shares are invalid |
| `config.toml` (kernel) | **UPDATE** trusted_height/hash | Must point to new chain's block |
| `story.toml` | **KEEP** | Unless kernel endpoints changed |

---

## DKG Stage Timing Reference

With current devnet params (registration=60, dealing=60, finalization=60, active=100):

| Stage | Duration | Cumulative (from round start) |
|-------|----------|-------------------------------|
| Registration | 60 blocks (~2 min) | 0-60 |
| Dealing | 60 blocks (~2 min) | 60-120 |
| Finalization | 60 blocks (~2 min) | 120-180 |
| Active | 100 blocks (~3.3 min) | 180-280 |

Total round duration: ~280 blocks (~9 minutes at 2s block time).

> These params are set in `client/app/upgrades/v_2_0_0/upgrades.go` → `dkgParamsForChain()` for `DKGTestChainID`.

---

## Troubleshooting

| Symptom | Cause | Fix |
|---------|-------|-----|
| App hash mismatch | Genesis or state inconsistency | `./story rollback` (1 block rollback), then restart |
| Kernel port 50051 still in use | Process didn't stop cleanly | `sudo kill -9 $(sudo lsof -t -i :50051)` |
| Different mrenclave across validators | Non-reproducible build | Verify same commit, `make clean && make build-with-cpp && make all-gramine` |
| `Connected to kernel endpoint` missing | Kernel not running or wrong port | Check kernel process, restart story after kernel starts |
| `chain id should not be empty` | Missing light client config | Add `[light_client]` section to kernel config.toml |
| `too many pages allocated` | SGX enclave memory too small | Increase `enclave_size` to 4G+ in manifest, re-sign |
| `story-kernel.manifest.sgx does not exist` | Wrong WorkingDirectory in systemd | Update systemd `WorkingDirectory` to story-kernel source dir |
| Bootnode consensus failure | Bootnode has no kernel, DKG events cause state divergence | Set `dkg.enable = false` on bootnode, or stop bootnode |
| `no whitelisted enclave type found for code commitment` | New CC not whitelisted | Run `whitelistEnclaveType` (Part 2, Step 4) |
| `DKG service already running; skipping` | Duplicate goroutine for same round | Expected dedup behavior; new round always preempts |
| Finalization signature mismatch | Enclave type mismatch (stale `k.enclaveType`) | Fixed in commit `c5a72139`; ensure latest story binary |
| `cast send` tx not being mined | Gas price too low for geth miner | Always use `--gas-price 17000000000 --legacy` |
| `createValidator` reverts | Gas estimation fails for precompile | Use explicit `--gas-limit 500000` |
| Kernel light client witness failure | Witness node is down or unreachable | Update `witness_addrs` to point to healthy validators only |
| `round regression at height N` | `priv_validator_state.json` has stale round | Reset: `echo '{"height":"0","round":0,"step":0}' > priv_validator_state.json` |
| `version of store dkg mismatch` | DKG binary used with genesis created by non-DKG binary | Use matching binary for genesis creation, or start with v1.5.3 and upgrade at v2.0.0 height |
| `execution_block_hash` mismatch | Hash computed with wrong geth version, or placed at wrong path | Must use same geth version as remote; must be at `app_state.evmengine.params.execution_block_hash` |
| `finalization signature address mismatch` | MRENCLAVE changed but genesis alloc not regenerated | Regenerate alloc with new `SGX_CODE_COMMITMENT` matching the current kernel binary |
| `cosmovisor` conflicts with `story.service` | Both try to run story | `sudo systemctl disable cosmovisor` |
| `header belongs to another chain` in kernel | Kernel config `chain_id` doesn't match actual chain | Update kernel config.toml `chain_id` to match CL genesis `chain_id` |
| Geth peers 0 after restart | `static-nodes.json` not auto-loaded | Add peers manually via `admin_addPeer` RPC |
| CL data not wiped (systemd restart race) | `story.service` Restart=always recreates data during wipe | Stop story and verify `pgrep -x story` returns empty before wiping data |
