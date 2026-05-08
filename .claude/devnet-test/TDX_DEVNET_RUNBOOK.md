# TDX + SGX Devnet Testing Runbook

> Operational runbook for the **SGX + TDX 2-validator** DKG test net.
> Read this AFTER `DEVNET_RUNBOOK.md` — that doc is the canonical source for
> the SGX-only 3-validator workflow; this doc is the **delta** for adding TDX.

This runbook supersedes DEVNET_RUNBOOK.md only for the differences listed
under "What changes vs the SGX-only runbook" below; everything else (DCAP
deploy, collateral upload, geth init, kernel light-client config) stays the
same.

---

## Purpose

Validate that the on-chain DKG protocol completes successfully when validator
membership is **mixed** between Intel SGX (Gramine, MRENCLAVE-based identity)
and Intel TDX (paravisor vTPM, MRTD‖RTMR-based identity), and that CDR
threshold decryption works across the heterogeneous committee.

Acceptance:

1. **DKG round 1**: both validators register; both attestations pass; round
   finalizes with a single group key.
2. **CDR**: a payload encrypted under the round-1 group key decrypts via two
   partial decryptions (one SGX, one TDX) with t=n=2.
3. **Resharing round 2**: triggered automatically at the end of the active
   period; new group key generated; CDR still works under the new key.

---

## What changes vs the SGX-only runbook

| Aspect | SGX-only (DEVNET_RUNBOOK) | SGX + TDX (this doc) |
|---|---|---|
| Validators | 3 SGX | 1 SGX + 1 TDX |
| DKG params | `n=3 t=67%` | `n=2 t=100%` |
| story branch | `dkg/hans-temp-test` | `dkg/hans-tdx-test` |
| story-kernel branch | `hans/temp-test` | `feat/tdx-support` |
| TDX validation hook | not deployed | deployed in genesis (enclaveType=2) |
| TDX kernel build | n/a | `make build-tdx` (configfs-tsm + vTPM) |
| Bootnode | dedicated `boot` host | reused — SGX validator carries P2P bootstrap |
| `enc-type` in story.toml | `1` (SGX) on every validator | `1` on SGX validator, `2` on TDX validator |
| Cloud SKU | `Standard_DC4s_v3` ×3 | `Standard_DC4s_v3` ×1 + `Standard_DC4es_v5` ×1 |

Anything else (chain IDs, EL/CL genesis layout, DCAP deploy steps,
collateral upload) is identical to DEVNET_RUNBOOK.md.

---

## Infrastructure

VMs provisioned in the Pip Labs Azure subscription on 2026-05-08.

| Role          | Hostname                  | Public IP        | Private IP | SKU                  | TEE |
|---------------|---------------------------|------------------|------------|----------------------|-----|
| SGX validator | `weu-dkg-tdx-test-sgx1`   | `52.157.106.205` | `10.10.0.5`| `Standard_DC4s_v3`   | SGX |
| TDX validator | `weu-dkg-tdx-test-tdx1`   | `20.126.100.96`  | `10.10.0.4`| `Standard_DC4es_v6`  | TDX |

> The SGX validator doubles as the bootnode (CometBFT seed + geth peer
> bootstrap) — no dedicated bootnode in this 2-validator topology. This works
> because `seeds`, `persistent_peers`, and geth `--bootnodes` can all point
> to the SGX validator.

> **TDX SKU note**: This is `DC4es_v6` (Intel Xeon 5 / Granite Rapids), not
> `DC4es_v5` originally planned. Our `TDXValidationHook` accepts both V4 and
> V5 quote layouts; offset constants for measurement fields are
> version-independent, so v6 hardware works without code changes.

### Hardware verification (one-shot, captured at provisioning time)

```
SGX VM:
  Ubuntu 24.04.4 LTS, kernel 6.17.0-1013-azure
  /dev/sgx_enclave + /dev/sgx_provision present
  CPU flags: sgx, sgx_lc

TDX VM:
  Ubuntu 24.04.4 LTS, kernel 6.17.0-1010-azure-fde (>= 6.7 ✓)
  /dev/tpm0 + /dev/tpmrm0 (vTPM via paravisor)
  /sys/kernel/config/tsm/report/ available (configfs-tsm)
  dmesg: "Memory Encryption Features active: Intel TDX"
  Initial PCR digests (sha256):
    PCR 7  = 0x3B20E022416FDF61D72E4DA32B4354781BE3DE0608116976D28FFDAD8C341D2A
    PCR 11 = 0xB75AC98943A8CA3CA80C5BA9C0392C840C56F92C4E87FC36E3FB749F5026E06A
```

### Azure resource group

Provisioned by Pip Labs infra (West Europe region — note the `weu-` hostname
prefix). Reuse / teardown via the company subscription, not the personal
provisioning script in this repo.

### SSH access

Use the existing devnet key. Add to `~/.ssh/config`:

```ssh-config
Host story-sgx
  HostName 52.157.106.205
  User ubuntu
  IdentityFile ~/.ssh/dkg_new.pem
  StrictHostKeyChecking no

Host story-tdx
  HostName 20.126.100.96
  User ubuntu
  IdentityFile ~/.ssh/dkg_new.pem
  StrictHostKeyChecking no
```

Bulk run:
```bash
for h in story-sgx story-tdx; do echo "=== $h ==="; ssh $h "<command>"; done
```

### Tear down

VMs live in the company Azure subscription; coordinate with Pip Labs infra
when the test program ends. Between sessions, use the standard `az vm
deallocate ...` to stop billing without losing state. The naver-account
provisioning script (`story-kernel/.claude/tdx-testnet/azure-setup.sh`) is
preserved for reference but is not what brought up these VMs.

### Key directories on remote

Same as `DEVNET_RUNBOOK.md` §"Key Directories":

| Path                              | Description                           |
|-----------------------------------|---------------------------------------|
| `~/.story/story/config/`          | CL config, genesis, validator keys    |
| `~/.story/story/data/`            | CL chain data                         |
| `~/.story/geth/data/`             | EL data                               |
| `/opt/story-kernel/`              | Kernel state, config, keys (Gramine SGX **and** TDX both mount here) |
| `~/config/genesis-geth.json`      | EL genesis (source of truth)          |
| `~/story/`                        | story source code                     |
| `~/story-kernel/`                 | story-kernel source code              |

> **TDX-specific**: although the TDX backend does not run inside Gramine, it
> still uses `/opt/story-kernel/` as its data root for parity with the SGX
> path. Don't move this on the TDX VM; the kernel binary expects it.

### Service names

Same as the SGX runbook: `story`, `node-geth`, `story-kernel`.

### Critical rules (TDX-specific delta)

These rules are **additions** on top of `DEVNET_RUNBOOK.md` §"Critical Rules".

**TDX kernel build & boot**

- TDX kernel binary is built with `make build-tdx`, **not** `make build-sgx`
  / `make all-gramine`. Single binary supports only one TEE backend.
- TDX kernel runs as a **plain process** (under systemd or `nohup`),
  **not** under `gramine-sgx`. There is no manifest, no signing step.
- TDX kernel requires:
  - Linux kernel ≥ 6.7 (Ubuntu 24.04 confidential VM image satisfies this).
    Verify: `uname -r`
  - `/sys/kernel/config/tsm/report/` available (configfs-tsm).
    Verify: `ls /sys/kernel/config/tsm/report/`
  - `/dev/tpmrm0` (preferred) or `/dev/tpm0` reachable.
    Verify: `ls /dev/tpm*`
  - PCR 7 + 11 non-zero. Verify:
    `sudo tpm2_pcrread sha256:7,11`

**TDX bootstrap mode (one-shot, before strict mode)**

- On the very first boot of a TDX validator, `enclave/tdx/providers.go`
  has `ExpectedDigest: nil` for `default-tpm-pcrs-7-11`.
- The kernel emits a WARN log on startup containing the empirically
  measured PCR digest in copy-pasteable hex.
- Operator workflow:
  1. Capture the digest from `journalctl -u story-kernel --no-pager | grep "bootstrap mode"`
  2. Edit `enclave/tdx/providers.go` on the TDX VM (or in a feature commit on
     `feat/tdx-support`) to populate `ExpectedDigest: []byte{0x..., ...}`
     with the captured value.
  3. `make build-tdx` and restart the kernel. From here on, any boot whose
     PCR state diverges from the populated entry will `log.Fatal` and refuse
     to start (strict mode).
- Bootstrap mode is per-VM. If you reprovision the TDX VM, expect to repeat.

**TDX code commitment & on-chain whitelist**

- Kernel-side `CodeCommitment` for TDX is the native concatenation
  `MRTD || RTMR0 || RTMR1 || RTMR2 || RTMR3` (240 bytes). The kernel does
  **not** hash this client-side.
- On-chain `expectedCodeCommitment` is `bytes32`, so `TDXValidationHook`
  computes `keccak256(MRTD || RTMR0..3)` from the raw quote and compares
  against the whitelisted bytes32.
- `TDX_CODE_COMMITMENT` constant in `contracts/script/GenerateAlloc.s.sol`
  is initially a placeholder. After the TDX kernel boots in strict mode and
  produces its first quote, capture the keccak digest and amend
  `GenerateAlloc.s.sol` (or call `whitelistEnclaveType` directly via
  governance). Without the correct digest, registration will revert with
  "TDXValidationHook: Code commitment does not match".

**story.toml `enc-type` setting**

- SGX validator: `enc-type = 1`
- TDX validator: `enc-type = 2`
- This must match the `enclaveType` the validator presents in its
  `register` tx (`bytes32(uint256(enc-type))`). A mismatch causes the
  on-chain DKG keeper to call the wrong validation hook → revert.

**TDX kernel single-binary mode**

- Unlike SGX which uses Gramine's manifest argv, the TDX kernel takes its
  config via the same `--home /opt/story-kernel` path. Write
  `/opt/story-kernel/config.toml` with the same shape as the SGX one.

---

## Branch management

The `dkg/hans-tdx-test` branch is the deployment branch. It tracks
`release/1.6` plus four cherry-picks:

```
3bd91932 feat: integrate TDXValidationHook into devnet alloc + n=2 DKG params
a77f3a22 feat: set USE_DEPLOYER_AS_OWNER=true for devnet testing
144fa6c8 feat: setup devnet (real SGXValidationHook + DCAP)
28891ff8 fix(dkg): recover from panic in decrypt worker goroutine
eaef9f7b feat(dkg): add TDXValidationHook for Intel TDX attestation
```

When release/1.6 gets new commits, rebase:

```bash
git checkout dkg/hans-tdx-test
git fetch origin
git rebase origin/release/1.6
git push origin dkg/hans-tdx-test --force
```

Devnet config tweaks (e.g., updated `TDX_CODE_COMMITMENT`) should be amended
into the most recent `feat: integrate TDXValidationHook ...` commit:

```bash
# edit GenerateAlloc.s.sol ...
git add -A
git commit --amend --no-edit --no-verify    # --no-verify because pre-existing
                                              # forge upgrade tests still fail
git push origin dkg/hans-tdx-test --force
```

> The `--no-verify` is the same compromise the parent
> `feat: set USE_DEPLOYER_AS_OWNER=true ...` made: `forge test` reports 25
> upgrade-related failures whenever USE_DEPLOYER_AS_OWNER=true. The relevant
> tests for our work (`test/dkg/*`) all pass.

---

## Workflow Quick Reference

The full workflow follows `DEVNET_RUNBOOK.md` Part 1 with these
substitutions:

1. **Step 1 (Stop processes)** — same. On TDX VM, `gramine-sgx` and `loader`
   are not present; only the plain `story-kernel` binary needs killing.
2. **Step 2 (Pull code)** —
   - `cd ~/story && git fetch origin && git checkout -f dkg/hans-tdx-test && git reset --hard origin/dkg/hans-tdx-test`
   - `cd ~/story-kernel && git fetch origin && git checkout -f feat/tdx-support && git reset --hard origin/feat/tdx-support`
3. **Step 3 (Build kernel)** —
   - SGX VM: `cd ~/story-kernel && make all-sgx` (build-sgx + Gramine sign)
   - TDX VM: `cd ~/story-kernel && make all-tdx` (build-tdx, no signing)
4. **Step 4 (Build story)** — same; `make build` then install binary.
5. **Step 5 (Generate alloc + EL/CL genesis)** —
   - On local machine, run `forge script script/GenerateAlloc.s.sol --tc GenerateAlloc -vvv --chain-id 1511`
   - Note both the SGX hook proxy address and the **new TDX hook proxy address**
     printed by `setTDXValidationHook`. Copy these into your operations log.
   - Run `add-alloc-extras.py` exactly as in DEVNET_RUNBOOK §5.2.
   - Compute execution_block_hash same as §5.4.
   - CL genesis (`genesis-node.json`):
     - `chain_id = "story-68931"`
     - `app_state.evmengine.params.execution_block_hash` = base64 of EL
       genesis hash
     - **Validator set** must contain exactly TWO entries: the SGX VM's
       priv_validator pubkey and the TDX VM's. Get them with
       `python3 -c "import json; print(json.load(open('PATH'))['pub_key'])"`
       run on each VM after `story init`.
     - `app_state.dkg.params.registration_period = 500` so DCAP deploy
       finishes before round 1 expires.
6. **Step 6 (Clean data)** — same on both VMs.
7. **Step 7 (geth init + start services)** — same; verify genesis hashes
   match across both VMs.
8. **Step 7.5 (DCAP + hooks)** —
   - Run `DeployAllDevnet.s.sol` exactly as in §7.5.1 from the SGX VM (it has
     more IP allocation in the alloc).
   - **Connect BOTH hooks**, not just SGX:
     ```
     # Same $ATTESTATION for both hooks
     SGX_HOOK="<from GenerateAlloc output>"
     TDX_HOOK="<from GenerateAlloc output>"

     for HOOK in $SGX_HOOK $TDX_HOOK; do
       cast send $HOOK 'setAutomataValidationAddr(address)' $ATTESTATION \
         --rpc-url $RPC --private-key $KEY --legacy --gas-price 30000000000
       cast send $HOOK 'setTcbEvaluationDataNumber(uint32)' 18 \
         --rpc-url $RPC --private-key $KEY --legacy --gas-price 30000000000
     done
     ```
   - Whitelist is already done in genesis (both `enclaveType=1` and
     `enclaveType=2`); verify with `isEnclaveTypeWhitelisted` for each.
9. **Step 8 (Start kernels)** —
   - SGX VM: `gramine-sgx story-kernel` under systemd, port 50051. Same as
     DEVNET_RUNBOOK.
   - TDX VM (first boot only — bootstrap):
     ```
     /opt/story-kernel/build/story-kernel start --home /opt/story-kernel \
         > /tmp/kernel.log 2>&1
     ```
     Watch for the WARN line containing "bootstrap mode" + the PCR digest
     hex. Capture, paste into `enclave/tdx/providers.go`, rebuild, restart.
   - TDX VM (strict mode): `sudo systemctl start story-kernel`.
10. **Step 9 (Verify DKG round)** — same. Both validators must register
    within the registration period; both must finalize.

---

## Validation acceptance log template

Fill this in as the test progresses (one row per milestone):

| Milestone                                     | Block | Notes |
|-----------------------------------------------|-------|-------|
| Both VMs provisioned, SSH OK                  |       |       |
| EL+CL genesis distributed, hashes match       |       |       |
| DCAP stack deployed, V3QV has code            |       |       |
| Both hooks connected to Automata + tcb=18    |       |       |
| TDX bootstrap digest captured                 |       |       |
| TDX kernel restarted in strict mode           |       |       |
| TDX_CODE_COMMITMENT amended + redeployed      |       |       |
| SGX validator registered (round 1)            |       |       |
| TDX validator registered (round 1)            |       |       |
| Round 1 dealing complete                      |       |       |
| Round 1 finalized (group pubkey published)    |       |       |
| CDR encrypt + decrypt OK (round 1 key)        |       |       |
| Resharing round 2 auto-triggered              |       |       |
| Round 2 finalized                             |       |       |
| CDR still works under round 2 key             |       |       |

---

## Troubleshooting

**TDX kernel boots but `runSelfCheck` fails with "TPM2_GetCapability failed"**
- vTPM device not yet ready at boot. Wait 5s and restart the service.
  (The kernel's `selfCheckOnce` makes the first TEE op pay this cost; later
  ops won't re-check.)

**TDX kernel boots, prints bootstrap WARN, then exits**
- Bootstrap mode is supposed to *proceed*, not exit. If it exits, look at
  earlier log lines for a `log.Fatal`. Common cause: `/dev/tdx_guest` or
  configfs-tsm path missing → kernel registers fail-closed quote provider.

**On-chain registration reverts with "TDXValidationHook: Code commitment does not match"**
- The whitelisted `TDX_CODE_COMMITMENT` does not match the kernel's
  `keccak256(MRTD || RTMR0..3)`. Re-derive the digest from the kernel's
  current quote and amend GenerateAlloc.s.sol or call
  `whitelistEnclaveType` again with the correct value.

**On-chain registration reverts with "TDXValidationHook: Not a TDX quote"**
- The kernel emitted an SGX-shaped quote (TEE type byte 0 ≠ 0x81). Confirm
  the kernel was built with `make build-tdx` and not the noop / sgx
  variant. `story-kernel version` should mention TDX.

**One validator finalizes, the other does not**
- t=2/n=2 means **both** must finalize. Check the slow validator's logs:
  `journalctl -u story --no-pager | grep -i "DKG\|error" | tail -50`
  and `journalctl -u story-kernel --no-pager | tail -50`.

**DKG round expires without registration completing**
- Default `registration_period` is too short for DCAP deploy + collateral
  upload. Per DEVNET_RUNBOOK §"DKG Registration Period", set to 500 in CL
  genesis.

---

## Pointers

- Provisioning script: `story-kernel/.claude/tdx-testnet/azure-setup.sh`
- TDX backend doc: `story-kernel/enclave/tdx/doc.go` (vTPM-in-TCB
  requirements, supportedProviders rationale)
- TDX validation hook: `story/contracts/src/protocol/TDXValidationHook.sol`
- TDX validation hook tests: `story/contracts/test/dkg/tdx_validation_hook.t.sol`
- DKG genesis alloc: `story/contracts/script/GenerateAlloc.s.sol`
  - `setSGXValidationHook()`, `setTDXValidationHook()`
- SGX-only canonical runbook: `story/.claude/devnet-test/DEVNET_RUNBOOK.md`
