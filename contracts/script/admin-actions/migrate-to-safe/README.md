# Migration of current multisig to a Safe

Each step corresponds to a specific script in this directory that generates the necessary transaction payloads for the timelock controllers.

Current Predeploys and their upgradeability are governed by a `TimelockController`, operated by 2 multisigs (Proposer + Security Council)

![Migration diagram 1](./images/1.migration.png)

In order to change that multisig, we need to follow these steps:

## 1. Deploy a new TimelockController
Deploy a new TimelockController governed by the new multisigs (still Proposer + Security Council multisigs, but different technology).
In case the new multisig doesn't work for some reason, the Timelock starts with the old multisig also holding the correspondent roles.
To prevent the need for this operation again, the new Timelock is going to have the proposers be also root admins (to grant roles), at least for some time.

Script: [1.DeployNewTimelock.s.sol](./1.DeployNewTimelock.s.sol)

![Migration diagram 2](./images/2.migration.png)

## 2. Transfer ownership of ProxyAdmins (proxy upgradeability) from old to new TimelockController
This process is split into multiple steps to reduce risk:

### 2.1 Transfer ownership of the first quarter of ProxyAdmins
The old TimelockController transfers ownership of the first quarter of the predeploy ProxyAdmins to the new timelock.

Script: [2.1.TransferOwnershipProxyAdmin1.s.sol](./2.1.TransferOwnershipProxyAdmin1.s.sol)

### 2.2 Transfer ownership of the second quarter of ProxyAdmins
The old TimelockController transfers ownership of the second quarter of predeploy ProxyAdmins to the new timelock.

Script: [2.2.TransferOwnershipProxyAdmin2.s.sol](./2.2.TransferOwnershipProxyAdmin2.s.sol)

### 2.3 Transfer ownership of the third quarter of ProxyAdmins
The old TimelockController transfers ownership of the third quarter of predeploy ProxyAdmins to the new timelock.

Script: [2.3.TransferOwnershipProxyAdmin3.s.sol](./2.3.TransferOwnershipProxyAdmin3.s.sol)

### 2.4 Transfer ownership of the fourth quarter of ProxyAdmins
The old TimelockController transfers ownership of the final quarter of predeploy ProxyAdmins to the new timelock.

Script: [2.4.TransferOwnershipProxyAdmin4.s.sol](./2.4.TransferOwnershipProxyAdmin4.s.sol)

## 3. Transfer ownership of Proxies in use from old to new TimelockController

Transfering these needs 2 steps because they are `Owneable2Step` instead of `Owneable`. New owner must create a transaction calling `acceptOwnership()` to finish the process.

This process is similarly split into multiple steps to reduce risk.

### 3.1 Transfer ownership of the UpgradesEntrypoint
The old TimelockController transfers ownership of the UpgradesEntrypoint proxy to the new timelock.

Script: [3.1.TransferOwnershipUpgradesEntrypoint.s.sol](./3.1.TransferOwnershipUpgradesEntrypoint.s.sol)

### 3.2 Accept ownership of the UpgradesEntrypoint
The new timelock accepts ownership of the UpgradesEntrypoint.

Script: [3.2.ReceiveOwnershipUpgradesEntryPoint.s.sol](./3.2.ReceiveOwnershipUpgradesEntryPoint.s.sol)

### 3.3 Transfer ownership of UBIPool
The old TimelockController transfers ownership of UBIPool to the new timelock.

Script: [3.3.TransferOwnershipUBIPool.s.sol](./3.3.TransferOwnershipUBIPool.s.sol)

### 3.4 Accept ownership of UBIPool
The new timelock accepts ownership of UBIPool.

Script: [3.4.ReceiveOwnershipUBIPool.s.sol](./3.4.ReceiveOwnershipUBIPool.s.sol)

### 3.5 Transfer ownership of IPTokenStaking
The old TimelockController transfers ownership of IPTokenStaking to the new timelock.

Script: [3.5.TransferOwnershipIPTokenStaking.s.sol](./3.5.TransferOwnershipIPTokenStaking.s.sol)

### 3.6 Accept ownership of IPTokenStaking
The new timelock accepts ownership of IPTokenStaking.

Script: [3.6.ReceiveOwnershipIPTokenStaking.s.sol](./3.6.ReceiveOwnershipIPTokenStaking.s.sol)

We are now here:

![Migration diagram 3](./images/3.migration.png)

## 4. Finalize the migration
After checking that everything works correctly, old multisigs renounce roles to complete the transition.

![Migration diagram 4](./images/4.migration.png)

Script: [4.RenounceGovernanceRoles.sol](./4.RenounceGovernanceRoles.sol)

## Automated Migration Execution

For convenience, an automated script is provided to run all migration steps sequentially:

### run-migration.sh

The `run-migration.sh` script automates the migration process by executing all ownership transfer scripts in the correct order. Note: Timelock deployment must be run separately before using this script.

#### Prerequisites

Before running this script, ensure the new timelock has been deployed:

```bash
# Deploy the new timelock first
forge script script/admin-actions/migrate-to-safe/1.DeployNewTimelock.s.sol:DeployNewTimelock --rpc-url <RPC_URL> --broadcast -vvv
```

#### Usage

```bash
# Dry run (no broadcasting)
./run-migration.sh <RPC_URL>

# Execute with broadcasting
./run-migration.sh <RPC_URL> --broadcast
```

#### Examples

```bash
# Test against local node
./run-migration.sh http://localhost:8545

# Execute on mainnet (be very careful!)
./run-migration.sh https://eth-mainnet.g.alchemy.com/v2/YOUR_API_KEY --broadcast

# Execute on testnet
./run-migration.sh https://eth-goerli.g.alchemy.com/v2/YOUR_API_KEY --broadcast
```

#### Features

- **Sequential execution**: Runs 11 migration scripts in the correct order (timelock deployment excluded)
- **Error handling**: Stops on any failure and reports the specific script that failed
- **Safety checks**: Requires explicit confirmation when broadcasting is enabled
- **Progress tracking**: Shows detailed progress and execution time
- **Dry run support**: Test the entire flow without broadcasting transactions
- **Environment validation**: Checks for required tools and environment variables

#### Environment Variables

Set any environment variables required by your scripts before running with `--broadcast`:

```bash
# Example environment variables that might be needed
export OLD_TIMELOCK_PROPOSER="0x..."
export OLD_TIMELOCK_EXECUTOR="0x..." 
export OLD_TIMELOCK_CANCELLER="0x..."
export SAFE_TIMELOCK_PROPOSER="0x..."
export SAFE_TIMELOCK_EXECUTOR="0x..."
export SAFE_TIMELOCK_CANCELLER="0x..."
# Add any other required environment variables for your specific setup
```

#### Script Execution Order

The automation script runs the following scripts in order:

**Note**: `1.DeployNewTimelock.s.sol` must be run separately before using this script.

1. `2.1.TransferOwnershipProxyAdmin1.s.sol`
2. `2.2.TransferOwnershipProxyAdmin2.s.sol`
3. `2.3.TransferOwnershipProxyAdmin3.s.sol`
4. `2.4.TransferOwnershipProxyAdmin4.s.sol`
5. `3.1.TransferOwnershipUpgradesEntrypoint.s.sol`
6. `3.2.ReceiveOwnershipUpgradesEntryPoint.s.sol`
7. `3.3.TransferOwnershipUBIPool.s.sol`
8. `3.4.ReceiveOwnershipUBIPool.s.sol`
9. `3.5.TransferOwnershipIPTokenStaking.s.sol`
10. `3.6.ReceiveOwnershipIPTokenStaking.s.sol`
11. `4.RenounceGovernanceRoles.s.sol`

#### Safety Considerations

⚠️ **Important**: When using `--broadcast`, the script will execute all migration steps on the actual network. Ensure you:

1. Test thoroughly on a testnet first
2. Verify all environment variables are correctly set
3. Have sufficient ETH for gas fees
4. Double-check the RPC URL is correct
5. Understand that this is an irreversible process once broadcast

The script includes confirmation prompts and warnings to help prevent accidental execution.
