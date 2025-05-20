/* solhint-disable no-console */
/* solhint-disable max-line-length */
// SPDX-License-Identifier: MIT

pragma solidity 0.8.23;

import { console2 } from "forge-std/Test.sol";
import { Test } from "test/utils/Test.sol";
import { stdJson } from "forge-std/StdJson.sol";

import { TimelockController } from "@openzeppelin/contracts/governance/TimelockController.sol";
import { Ownable } from "@openzeppelin/contracts/access/Ownable.sol";

import { Predeploys } from "src/libraries/Predeploys.sol";
import { Create3 } from "src/deploy/Create3.sol";
import { JSONTxWriter } from "script/utils/JSONTxWriter.s.sol";
import { EIP1967Helper } from "script/utils/EIP1967Helper.sol";
// Import migration scripts
import { DeployNewTimelock } from "script/admin-actions/migrate-to-safe/1.DeployNewTimelock.s.sol";
import { BaseTransferOwnershipProxyAdmin } from "script/admin-actions/migrate-to-safe/BaseTransferOwnershipProxyAdmin.s.sol";
import { TransferOwnershipsProxyAdmin1 } from "script/admin-actions/migrate-to-safe/2.1.TransferOwnershipProxyAdmin1.s.sol";
import { TransferOwnershipsProxyAdmin2 } from "script/admin-actions/migrate-to-safe/2.2.TransferOwnershipProxyAdmin2.s.sol";
import { TransferOwnershipsProxyAdmin3 } from "script/admin-actions/migrate-to-safe/2.3.TransferOwnershipProxyAdmin3.s.sol";
import { TransferOwnershipsProxyAdmin4 } from "script/admin-actions/migrate-to-safe/2.4.TransferOwnershipProxyAdmin4.s.sol";
import { TransferOwnershipsUpgradesEntrypoint } from "script/admin-actions/migrate-to-safe/3.1.TransferOwnershipUpgradesEntrypoint.s.sol";
import { ReceiveOwnershipUpgradesEntryPoint } from "script/admin-actions/migrate-to-safe/3.2.ReceiveOwnershipUpgradesEntryPoint.s.sol";
import { TransferOwnershipsRestPredeploys } from "script/admin-actions/migrate-to-safe/3.3.TransferOwnershipRestPredeploys.s.sol";
import { ReceiveOwnershipRestPredeploys } from "script/admin-actions/migrate-to-safe/3.4.ReceiveOwnershipRestPredeploys.s.sol";

contract SafeMigrationCancelScriptsTest is Test {
    using stdJson for string;

    // Mock addresses for the test
    address private OLD_TIMELOCK_PROPOSER;
    address private OLD_TIMELOCK_GUARDIAN;
    address private OLD_TIMELOCK_EXECUTOR;
    address private constant SAFE_TIMELOCK_PROPOSER = address(0x1111111111111111111111111111111111111111);
    address private constant SAFE_TIMELOCK_EXECUTOR = address(0x2222222222222222222222222222222222222222);
    address private constant SAFE_TIMELOCK_GUARDIAN = address(0x3333333333333333333333333333333333333333);

    // Private key for timelock deployer
    uint256 private constant DEPLOYER_PRIVATE_KEY = 0x1;

    // Maximum number of transactions in JSON files
    uint256 private constant MAX_TXS_PER_JSON = 1000;

    // Timelock controller addresses
    address private oldTimelockAddress;
    address private newTimelockAddress;
    TimelockController private oldTimelock;
    TimelockController private newTimelock;

    // Mock contracts for testing
    address[] private proxyAdmins;
    // JSON directories
    string constant TIMELOCK_ROOT_DIR = "./broadcast/";
    string constant OUTPUT_DIR = "./script/admin-actions/output/";

    // Role constants from TimelockController
    bytes32 private constant DEFAULT_ADMIN_ROLE = 0x00;

    uint256 private constant MIN_DELAY = 300;

    function setUp() public override {
        super.setUp();
        OLD_TIMELOCK_PROPOSER = admin;
        OLD_TIMELOCK_EXECUTOR = admin;
        OLD_TIMELOCK_GUARDIAN = admin;

        oldTimelock = timelock;

        // Setup environment variables
        _setupEnvVars();
    }

    function testCancelMigrationScripts() public {
        // Step 1: Deploy new timelock
        _testDeployNewTimelock();

        // Step 2: Transfer ownership of ProxyAdmins
        _testCancelTransferOwnershipProxyAdmins(1);
        _testCancelTransferOwnershipProxyAdmins(2);
        _testCancelTransferOwnershipProxyAdmins(3);
        _testCancelTransferOwnershipProxyAdmins(4);

        // Step 3.1: Transfer ownership of UpgradesEntrypoint
        _testCancelTransferOwnershipUpgradesEntrypoint();

        // Step 3.2: Accept ownership of UpgradesEntrypoint
        _testCancelReceiveOwnershipUpgradesEntrypoint();

        // Step 3.3: Transfer ownership of rest of predeploys
        _testCancelTransferOwnershipsRestPredeploys();

        // Step 3.4: Accept ownership of rest of predeploys
        _testCancelReceiveOwnershipRestPredeploys();
    }

    function _setupEnvVars() private {
        // Setup environment variables
        vm.setEnv("NEW_TIMELOCK_DEPLOYER_PRIVATE_KEY", vm.toString(DEPLOYER_PRIVATE_KEY));
        vm.setEnv("SAFE_TIMELOCK_PROPOSER", vm.toString(SAFE_TIMELOCK_PROPOSER));
        vm.label(SAFE_TIMELOCK_PROPOSER, "SafeTimelockProposer");
        vm.setEnv("SAFE_TIMELOCK_EXECUTOR", vm.toString(SAFE_TIMELOCK_EXECUTOR));
        vm.label(SAFE_TIMELOCK_EXECUTOR, "SafeTimelockExecutor");
        vm.setEnv("SAFE_TIMELOCK_GUARDIAN", vm.toString(SAFE_TIMELOCK_GUARDIAN));
        vm.label(SAFE_TIMELOCK_GUARDIAN, "SafeTimelockGuardian");
        vm.setEnv("OLD_TIMELOCK_ADDRESS", vm.toString(address(oldTimelock)));
        vm.label(address(oldTimelock), "OldTimelock");
        vm.setEnv("OLD_TIMELOCK_PROPOSER", vm.toString(OLD_TIMELOCK_PROPOSER));
        vm.label(OLD_TIMELOCK_PROPOSER, "Old Timelock Proposer");
        vm.setEnv("OLD_TIMELOCK_EXECUTOR", vm.toString(OLD_TIMELOCK_EXECUTOR));
        vm.label(OLD_TIMELOCK_EXECUTOR, "Old Timelock Executor");
        vm.setEnv("OLD_TIMELOCK_GUARDIAN", vm.toString(OLD_TIMELOCK_GUARDIAN));
        vm.label(OLD_TIMELOCK_GUARDIAN, "Old Timelock Guardian");
        vm.setEnv("MIN_DELAY", vm.toString(MIN_DELAY)); // 5 minutes
    }

    function _testDeployNewTimelock() private {
        DeployNewTimelock deployScript = new DeployNewTimelock();
        deployScript.run();

        // Get new timelock address
        bytes32 salt = keccak256("STORY_TIMELOCK_CONTROLLER_SAFE");
        newTimelockAddress = Create3(Predeploys.Create3).predictDeterministicAddress(salt);
        newTimelock = TimelockController(payable(newTimelockAddress));

        // Verify new timelock setup
        assertTrue(newTimelockAddress.code.length > 0, "New timelock not deployed");
        assertTrue(newTimelock.hasRole(newTimelock.PROPOSER_ROLE(), SAFE_TIMELOCK_PROPOSER), "Safe admin not proposer");
        assertTrue(
            newTimelock.hasRole(newTimelock.PROPOSER_ROLE(), OLD_TIMELOCK_PROPOSER),
            "Old multisig not proposer"
        );
        assertTrue(
            newTimelock.hasRole(newTimelock.EXECUTOR_ROLE(), SAFE_TIMELOCK_EXECUTOR),
            "Safe executor not executor"
        );
        assertTrue(
            newTimelock.hasRole(newTimelock.EXECUTOR_ROLE(), OLD_TIMELOCK_EXECUTOR),
            "Old multisig not executor"
        );
        assertTrue(
            newTimelock.hasRole(newTimelock.CANCELLER_ROLE(), SAFE_TIMELOCK_GUARDIAN),
            "Safe guardian not canceller"
        );
        assertTrue(
            newTimelock.hasRole(newTimelock.CANCELLER_ROLE(), OLD_TIMELOCK_GUARDIAN),
            "Old security council not canceller"
        );
        assertTrue(newTimelock.hasRole(DEFAULT_ADMIN_ROLE, SAFE_TIMELOCK_PROPOSER), "New Safe admin is not admin");
        assertTrue(newTimelock.hasRole(DEFAULT_ADMIN_ROLE, OLD_TIMELOCK_PROPOSER), "Old multisig admin is not admin");
        assertFalse(newTimelock.hasRole(DEFAULT_ADMIN_ROLE, deployer), "deployer is admin");
        assertEq(newTimelock.getMinDelay(), MIN_DELAY, "Delay not set correctly");
    }

    /**
     * @notice Common functionality to schedule and cancel a timelock operation
     * @param timelockController The timelock controller to use (oldTimelock or newTimelock)
     * @param jsonFilename The base filename to read transactions from
     */
    function _scheduleAndCancelOperation(TimelockController timelockController, string memory jsonFilename) private {
        // Get all transaction JSONs (schedule, cancel, execute)
        (
            JSONTxWriter.Transaction memory scheduleTx,
            JSONTxWriter.Transaction memory executeTx,
            JSONTxWriter.Transaction memory cancelTx
        ) = _readAllTransactionFiles(jsonFilename);

        // Execute the schedule transaction
        console2.log("Scheduling transaction");
        _rawTimelockTransaction(scheduleTx);

        // Extract operation ID directly from the cancel transaction data
        bytes32 operationId = _getOperationIdFromCancelData(cancelTx.data);

        // Check that the operation is scheduled
        assertTrue(timelockController.isOperationPending(operationId), "Operation not scheduled");

        // Execute the cancel transaction
        console2.log("Canceling transaction");
        _rawTimelockTransaction(cancelTx);

        // Verify that the operation is no longer pending after cancellation
        assertFalse(timelockController.isOperationPending(operationId), "Operation still pending after cancel");

        // Wait for the timelock delay
        console2.log("Waiting for timelock delay");
        vm.warp(block.timestamp + MIN_DELAY + 1);

        // Attempt to execute should fail
        vm.startPrank(executeTx.from);
        vm.expectRevert("TimelockController: operation is not ready");
        (bool success, ) = executeTx.to.call{ value: executeTx.value }(executeTx.data);
        assertFalse(success, "Execute transaction should fail");
        vm.stopPrank();
    }

    function _testCancelTransferOwnershipProxyAdmins(uint256 iteration) private {
        BaseTransferOwnershipProxyAdmin script;
        if (iteration == 1) {
            script = BaseTransferOwnershipProxyAdmin(new TransferOwnershipsProxyAdmin1());
        } else if (iteration == 2) {
            script = BaseTransferOwnershipProxyAdmin(new TransferOwnershipsProxyAdmin2());
        } else if (iteration == 3) {
            script = BaseTransferOwnershipProxyAdmin(new TransferOwnershipsProxyAdmin3());
        } else if (iteration == 4) {
            script = BaseTransferOwnershipProxyAdmin(new TransferOwnershipsProxyAdmin4());
        } else {
            revert("invalid iteration");
        }

        // Run the script to generate JSON
        script.run();

        // Schedule and cancel the operation
        _scheduleAndCancelOperation(oldTimelock, script.message());

        // Verify that ownership did not transfer
        for (uint160 i = uint160(script.fromIndex()); i <= script.toIndex(); i++) {
            // Calculate proxy address using namespace offset
            address proxyAddress = address(uint160(uint160(Predeploys.Namespace) + i));

            // Get proxy admin address from proxy using EIP1967 storage slot
            Ownable proxyAdmin = Ownable(EIP1967Helper.getAdmin(proxyAddress));

            // Verify proxy admin is still owned by old timelock
            assertEq(proxyAdmin.owner(), address(oldTimelock), "Proxy admin transferred despite cancellation");
        }
    }

    function _testCancelTransferOwnershipUpgradesEntrypoint() private {
        // Run the script to generate JSON
        TransferOwnershipsUpgradesEntrypoint script = new TransferOwnershipsUpgradesEntrypoint();
        script.run();

        // Schedule and cancel the operation
        _scheduleAndCancelOperation(oldTimelock, "safe-migr-transfer-ownerships-upgrades-entrypoint");

        // Verify that the UpgradesEntrypoint has not transferred ownership
        assertEq(
            upgradeEntrypoint.pendingOwner(),
            address(0),
            "UpgradesEntrypoint ownership transferred despite cancellation"
        );
    }

    function _testCancelReceiveOwnershipUpgradesEntrypoint() private {
        // For this test, we need to simulate that ownership was first transferred
        // Set the pending owner of UpgradesEntrypoint to the new timelock
        vm.startPrank(address(oldTimelock));
        upgradeEntrypoint.transferOwnership(newTimelockAddress);
        vm.stopPrank();

        // Run the script to generate JSON
        ReceiveOwnershipUpgradesEntryPoint script = new ReceiveOwnershipUpgradesEntryPoint();
        script.run();

        // Schedule and cancel the operation
        _scheduleAndCancelOperation(newTimelock, "safe-migr-receive-ownerships-upgrades-entrypoint");

        // Verify that the UpgradesEntrypoint has not claimed ownership
        assertEq(
            upgradeEntrypoint.owner(),
            address(oldTimelock),
            "UpgradesEntrypoint ownership claimed despite cancellation"
        );
        assertEq(upgradeEntrypoint.pendingOwner(), newTimelockAddress, "UpgradesEntrypoint pending owner changed");
    }

    function _testCancelTransferOwnershipsRestPredeploys() private {
        // Run the script to generate JSON
        TransferOwnershipsRestPredeploys script = new TransferOwnershipsRestPredeploys();
        script.run();

        // Schedule and cancel the operation
        _scheduleAndCancelOperation(oldTimelock, "safe-migr-transfer-ownerships-rest-predeploys");

        // Verify that UBIPool and IPTokenStaking do not have new pending owner
        assertEq(ubiPool.pendingOwner(), address(0), "UBIPool ownership transferred despite cancellation");
        assertEq(
            ipTokenStaking.pendingOwner(),
            address(0),
            "IPTokenStaking ownership transferred despite cancellation"
        );
    }

    function _testCancelReceiveOwnershipRestPredeploys() private {
        // For this test, we need to simulate that ownership was first transferred
        vm.startPrank(address(oldTimelock));
        ubiPool.transferOwnership(newTimelockAddress);
        ipTokenStaking.transferOwnership(newTimelockAddress);
        vm.stopPrank();

        // Run the script to generate JSON
        ReceiveOwnershipRestPredeploys script = new ReceiveOwnershipRestPredeploys();
        script.run();

        // Schedule and cancel the operation
        _scheduleAndCancelOperation(newTimelock, "safe-migr-receive-ownerships-rest-predeploys");

        // Verify that the rest of the predeploys have not claimed new ownership
        assertEq(ubiPool.owner(), address(oldTimelock), "UBIPool ownership changed despite cancellation");
        assertEq(ipTokenStaking.owner(), address(oldTimelock), "IPTokenStaking ownership changed despite cancellation");
        assertEq(ubiPool.pendingOwner(), newTimelockAddress, "UBIPool pending owner changed");
        assertEq(ipTokenStaking.pendingOwner(), newTimelockAddress, "IPTokenStaking pending owner changed");
    }

    /**
     * @notice Extract operation ID from timelock cancel function call data
     * @param data The call data from a cancel transaction
     * @return operationId The extracted operation ID
     */
    function _getOperationIdFromCancelData(bytes memory data) internal pure returns (bytes32 operationId) {
        // The cancel function has the signature: cancel(bytes32)
        // The first 4 bytes are the function selector, and the next 32 bytes are the operationId
        require(
            bytes4(abi.encodePacked(data[0], data[1], data[2], data[3])) == bytes4(keccak256("cancel(bytes32)")),
            "Not a cancel function call"
        );

        bytes memory operationIdBytes = new bytes(32);
        for (uint i = 0; i < 32; i++) {
            operationIdBytes[i] = data[i + 4];
        }

        assembly {
            operationId := mload(add(operationIdBytes, 32))
        }

        return operationId;
    }

    /**
     * @notice Read transactions from schedule, cancel, and execute JSON files
     * @param baseFilename The base filename without suffix (-schedule, -cancel, -execute)
     * @return scheduleTx Transaction struct from schedule file
     * @return executeTx Transaction struct from execute file
     * @return cancelTx Transaction struct from cancel file
     */
    function _readAllTransactionFiles(
        string memory baseFilename
    )
        internal
        returns (
            JSONTxWriter.Transaction memory scheduleTx,
            JSONTxWriter.Transaction memory executeTx,
            JSONTxWriter.Transaction memory cancelTx
        )
    {
        // Create paths for all three file types
        string memory basePath = string.concat(OUTPUT_DIR, vm.toString(block.chainid), "/");
        string memory schedulePath = string.concat(basePath, baseFilename, "-schedule.json");
        string memory cancelPath = string.concat(basePath, baseFilename, "-cancel.json");
        string memory executePath = string.concat(basePath, baseFilename, "-execute.json");

        // Read schedule transaction
        assertTrue(vm.exists(schedulePath), "Schedule JSON file not found");
        string memory scheduleJson = vm.readFile(schedulePath);
        JSONTxWriter.Transaction[] memory scheduleTxs = _parseTransactionsFromJson(scheduleJson);
        assertEq(scheduleTxs.length, 1, "Schedule JSON must contain exactly one transaction");
        scheduleTx = scheduleTxs[0];

        // Read cancel transaction
        assertTrue(vm.exists(cancelPath), "Cancel JSON file not found");
        string memory cancelJson = vm.readFile(cancelPath);
        JSONTxWriter.Transaction[] memory cancelTxs = _parseTransactionsFromJson(cancelJson);
        assertEq(cancelTxs.length, 1, "Cancel JSON must contain exactly one transaction");
        cancelTx = cancelTxs[0];

        // Read execute transaction
        assertTrue(vm.exists(executePath), "Execute JSON file not found");
        string memory executeJson = vm.readFile(executePath);
        JSONTxWriter.Transaction[] memory executeTxs = _parseTransactionsFromJson(executeJson);
        assertEq(executeTxs.length, 1, "Execute JSON must contain exactly one transaction");
        executeTx = executeTxs[0];

        return (scheduleTx, executeTx, cancelTx);
    }

    function _readRegularTransactionFiles(
        string memory baseFilename
    ) internal returns (JSONTxWriter.Transaction[] memory txs) {
        string memory basePath = string.concat(OUTPUT_DIR, vm.toString(block.chainid), "/");
        string memory path = string.concat(basePath, baseFilename, "-regular.json");
        assertTrue(vm.exists(path), "Regular JSON file not found");
        string memory json = vm.readFile(path);
        txs = _parseTransactionsFromJson(json);
    }

    /**
     * @notice Parse a JSON string into an array of Transaction structs
     * @param json The JSON string to parse
     * @return An array of Transaction structs
     */
    function _parseTransactionsFromJson(string memory json) internal view returns (JSONTxWriter.Transaction[] memory) {
        // Get the number of transactions in the JSON array
        // Create an array to store the transactions
        JSONTxWriter.Transaction[] memory readTxs = new JSONTxWriter.Transaction[](MAX_TXS_PER_JSON);
        uint256 effectiveTxs = 0;
        // Parse each transaction in the array
        for (uint256 i = 0; i < MAX_TXS_PER_JSON; i++) {
            try this._parseTransaction(json, i) returns (JSONTxWriter.Transaction memory transaction) {
                readTxs[i] = transaction;
                effectiveTxs++;
            } catch {
                console2.log("No more transactions in JSON");
                break;
            }
        }

        JSONTxWriter.Transaction[] memory transactions = new JSONTxWriter.Transaction[](effectiveTxs);
        for (uint256 i = 0; i < effectiveTxs; i++) {
            transactions[i] = readTxs[i];
        }
        return transactions;
    }

    function _parseTransaction(
        string memory json,
        uint256 index
    ) external pure returns (JSONTxWriter.Transaction memory transaction) {
        string memory indexPath = string.concat("[", vm.toString(index), "]");

        address from = stdJson.readAddress(json, string.concat(indexPath, ".from"));
        address to = stdJson.readAddress(json, string.concat(indexPath, ".to"));
        uint256 value = stdJson.readUint(json, string.concat(indexPath, ".value"));
        bytes memory data = stdJson.readBytes(json, string.concat(indexPath, ".data"));
        uint8 operation = uint8(stdJson.readUint(json, string.concat(indexPath, ".operation")));
        string memory comment = stdJson.readString(json, string.concat(indexPath, ".comment"));

        // Create the transaction struct
        return
            JSONTxWriter.Transaction({
                from: from,
                to: to,
                value: value,
                data: data,
                operation: operation,
                comment: comment
            });
    }

    /**
     * @notice Execute a single transaction from the old timelock
     * @param transaction The transaction to execute
     */
    function _rawTimelockTransaction(JSONTxWriter.Transaction memory transaction) internal {
        vm.startPrank(transaction.from);
        (bool success, ) = transaction.to.call{ value: transaction.value }(transaction.data);
        require(success, "Transaction execution failed");
        vm.stopPrank();
    }
}
