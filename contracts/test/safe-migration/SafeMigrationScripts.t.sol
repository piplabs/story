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
import { RenounceGovernanceRoles } from "script/admin-actions/migrate-to-safe/4.RenounceGovernanceRoles.s.sol";

contract SafeMigrationScriptsTest is Test {
    using stdJson for string;

    // Mock addresses for the test
    address private OLD_TIMELOCK_PROPOSER;
    address private OLD_TIMELOCK_CANCELLER;
    address private OLD_TIMELOCK_EXECUTOR;
    address private constant SAFE_TIMELOCK_PROPOSER = address(0x1111111111111111111111111111111111111111);
    address private constant SAFE_TIMELOCK_EXECUTOR = address(0x2222222222222222222222222222222222222222);
    address private constant SAFE_TIMELOCK_CANCELLER = address(0x3333333333333333333333333333333333333333);

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
        // Create output directory if it doesn't exist
        string memory outputDir = string.concat(OUTPUT_DIR, "31337/");
        vm.createDir(outputDir, true);

        OLD_TIMELOCK_PROPOSER = admin;
        OLD_TIMELOCK_EXECUTOR = admin;
        OLD_TIMELOCK_CANCELLER = admin;

        oldTimelock = timelock;

        // Setup environment variables
        _setupEnvVars();
    }

    function testMigrationScripts() public {
        // Step 1: Deploy new timelock
        _testDeployNewTimelock();

        // Step 2: Transfer ownership of ProxyAdmins
        _testTransferOwnershipProxyAdmins(1);
        _testTransferOwnershipProxyAdmins(2);
        _testTransferOwnershipProxyAdmins(3);
        _testTransferOwnershipProxyAdmins(4);

        // Step 3.1: Transfer ownership of UpgradesEntrypoint
        _testTransferOwnershipUpgradesEntrypoint();

        // Step 3.2: Accept ownership of UpgradesEntrypoint
        _testReceiveOwnershipUpgradesEntrypoint();

        // Step 3.3: Transfer ownership of rest of predeploys
        _testTransferOwnershipsRestPredeploys();

        // Step 3.4: Accept ownership of rest of predeploys
        _testReceiveOwnershipRestPredeploys();

        // Step 4: Renounce ownership of old multisig
        _testRenounceOwnershipOldMultisig();
    }

    function _setupEnvVars() private {
        // Setup environment variables
        vm.setEnv("NEW_TIMELOCK_DEPLOYER_PRIVATE_KEY", vm.toString(DEPLOYER_PRIVATE_KEY));
        vm.setEnv("SAFE_TIMELOCK_PROPOSER", vm.toString(SAFE_TIMELOCK_PROPOSER));
        vm.label(SAFE_TIMELOCK_PROPOSER, "SafeTimelockProposer");
        vm.setEnv("SAFE_TIMELOCK_EXECUTOR", vm.toString(SAFE_TIMELOCK_EXECUTOR));
        vm.label(SAFE_TIMELOCK_EXECUTOR, "SafeTimelockExecutor");
        vm.setEnv("SAFE_TIMELOCK_CANCELLER", vm.toString(SAFE_TIMELOCK_CANCELLER));
        vm.label(SAFE_TIMELOCK_CANCELLER, "SafeTimelockGuardian");
        vm.setEnv("OLD_TIMELOCK_ADDRESS", vm.toString(address(oldTimelock)));
        vm.label(address(oldTimelock), "OldTimelock");
        vm.setEnv("OLD_TIMELOCK_PROPOSER", vm.toString(OLD_TIMELOCK_PROPOSER));
        vm.label(OLD_TIMELOCK_PROPOSER, "Old Timelock Proposer");
        vm.setEnv("OLD_TIMELOCK_EXECUTOR", vm.toString(OLD_TIMELOCK_EXECUTOR));
        vm.label(OLD_TIMELOCK_EXECUTOR, "Old Timelock Executor");
        vm.setEnv("OLD_TIMELOCK_CANCELLER", vm.toString(OLD_TIMELOCK_CANCELLER));
        vm.label(OLD_TIMELOCK_CANCELLER, "Old Timelock Guardian");
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
            newTimelock.hasRole(newTimelock.CANCELLER_ROLE(), SAFE_TIMELOCK_CANCELLER),
            "Safe guardian not canceller"
        );
        assertTrue(
            newTimelock.hasRole(newTimelock.CANCELLER_ROLE(), OLD_TIMELOCK_CANCELLER),
            "Old security council not canceller"
        );
        assertTrue(newTimelock.hasRole(DEFAULT_ADMIN_ROLE, SAFE_TIMELOCK_PROPOSER), "New Safe admin is not admin");
        assertTrue(newTimelock.hasRole(DEFAULT_ADMIN_ROLE, OLD_TIMELOCK_PROPOSER), "Old multisig admin is not admin");
        assertFalse(newTimelock.hasRole(DEFAULT_ADMIN_ROLE, deployer), "deployer is admin");
        assertEq(newTimelock.getMinDelay(), MIN_DELAY, "Delay not set correctly");
    }

    function _testTransferOwnershipProxyAdmins(uint256 iteration) private {
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

        // Get all transaction JSONs (schedule, cancel, execute)
        (
            JSONTxWriter.Transaction memory scheduleTx,
            JSONTxWriter.Transaction memory executeTx,
            JSONTxWriter.Transaction memory cancelTx
        ) = _readAllTransactionFiles(script.message());

        // Execute the full timelock flow with verification
        _rawTimelockTransaction(scheduleTx);

        // Wait for the timelock delay
        console2.log("Waiting for timelock delay");
        vm.warp(block.timestamp + MIN_DELAY + 1);

        // Execute execute transaction
        console2.log("Executing execute transaction");
        _rawTimelockTransaction(executeTx);

        for (uint160 i = uint160(script.fromIndex()); i <= script.toIndex(); i++) {
            // Calculate proxy address using namespace offset
            address proxyAddress = address(uint160(uint160(Predeploys.Namespace) + i));

            // Get proxy admin address from proxy using EIP1967 storage slot
            Ownable proxyAdmin = Ownable(EIP1967Helper.getAdmin(proxyAddress));

            // Verify proxy admin is now the new timelock
            assertEq(proxyAdmin.owner(), newTimelockAddress, "Proxy admin not transferred");
        }
    }

    function _testTransferOwnershipProxyAdmin2() private {
        // Run the script to generate JSON
        TransferOwnershipsProxyAdmin2 script = new TransferOwnershipsProxyAdmin2();
        script.run(); // Generate operation JSON

        // Get all transaction JSONs (schedule, cancel, execute)
        (
            JSONTxWriter.Transaction memory scheduleTx,
            JSONTxWriter.Transaction memory executeTx,
            JSONTxWriter.Transaction memory cancelTx
        ) = _readAllTransactionFiles("safe-migr-transfer-ownerships-proxy-admin-2");

        // Execute the full timelock flow with verification
        _rawTimelockTransaction(scheduleTx);

        // Wait for the timelock delay
        console2.log("Waiting for timelock delay");
        vm.warp(block.timestamp + MIN_DELAY + 1);

        // Execute execute transaction
        console2.log("Executing execute transaction");
        _rawTimelockTransaction(executeTx);

        // Verify that half of the ProxyAdmins have new ownership
        uint256 halfSize = Predeploys.NamespaceSize / 2;
        uint256 startIdx = 1; // First half of proxies

        for (uint160 i = uint160(startIdx); i < halfSize; i++) {
            // Calculate proxy address using namespace offset
            address proxyAddress = address(uint160(uint160(Predeploys.Namespace) + i));

            // Get proxy admin address from proxy using EIP1967 storage slot
            Ownable proxyAdmin = Ownable(EIP1967Helper.getAdmin(proxyAddress));

            // Verify proxy admin is now the new timelock
            assertEq(proxyAdmin.owner(), newTimelockAddress, "Proxy admin not transferred");
        }
    }

    function _testTransferOwnershipUpgradesEntrypoint() private {
        // Run the script to generate JSON
        TransferOwnershipsUpgradesEntrypoint script = new TransferOwnershipsUpgradesEntrypoint();
        script.run(); // Generate operation JSON

        // Get all transaction JSONs (schedule, cancel, execute)
        (
            JSONTxWriter.Transaction memory scheduleTx,
            JSONTxWriter.Transaction memory executeTx,
            JSONTxWriter.Transaction memory cancelTx
        ) = _readAllTransactionFiles("safe-migr-transfer-ownerships-upgrades-entrypoint");

        // Execute the full timelock flow with verification
        _rawTimelockTransaction(scheduleTx);

        // Wait for the timelock delay
        console2.log("Waiting for timelock delay");
        vm.warp(block.timestamp + MIN_DELAY + 1);

        // Execute execute transaction
        console2.log("Executing execute transaction");
        _rawTimelockTransaction(executeTx);

        // Verify that the UpgradesEntrypoint has new ownership pending
        assertEq(upgradeEntrypoint.pendingOwner(), newTimelockAddress, "UpgradesEntrypoint not transferred");
    }

    function _testReceiveOwnershipUpgradesEntrypoint() private {
        // Run the script to generate JSON
        ReceiveOwnershipUpgradesEntryPoint script = new ReceiveOwnershipUpgradesEntryPoint();
        script.run(); // Generate operation JSON

        // Get all transaction JSONs (schedule, cancel, execute)
        (
            JSONTxWriter.Transaction memory scheduleTx,
            JSONTxWriter.Transaction memory executeTx,
            JSONTxWriter.Transaction memory cancelTx
        ) = _readAllTransactionFiles("safe-migr-receive-ownerships-upgrades-entrypoint");

        // Execute the full timelock flow with verification
        require(scheduleTx.to == newTimelockAddress, "Execute transaction is not to the new timelock");
        require(
            newTimelock.hasRole(newTimelock.PROPOSER_ROLE(), scheduleTx.from),
            "Schedule transaction is not from the new timelock"
        );
        _rawTimelockTransaction(scheduleTx);

        // Wait for the timelock delay
        console2.log("Waiting for timelock delay");
        vm.warp(block.timestamp + MIN_DELAY + 1);

        // Execute execute transaction
        console2.log("Executing execute transaction");
        require(executeTx.to == newTimelockAddress, "Execute transaction is not to the new timelock");
        require(
            newTimelock.hasRole(newTimelock.EXECUTOR_ROLE(), executeTx.from),
            "Execute transaction is not from the new timelock"
        );
        _rawTimelockTransaction(executeTx);

        // Verify that the UpgradesEntrypoint has new ownership
        assertEq(upgradeEntrypoint.owner(), newTimelockAddress, "UpgradesEntrypoint not transferred");
        assertEq(upgradeEntrypoint.pendingOwner(), address(0), "UpgradesEntrypoint pending owner not cleared");
    }

    function _testTransferOwnershipsRestPredeploys() private {
        // Run the script to generate JSON
        TransferOwnershipsRestPredeploys script = new TransferOwnershipsRestPredeploys();
        script.run(); // Generate operation JSON

        // Get all transaction JSONs (schedule, cancel, execute)
        (
            JSONTxWriter.Transaction memory scheduleTx,
            JSONTxWriter.Transaction memory executeTx,
            JSONTxWriter.Transaction memory cancelTx
        ) = _readAllTransactionFiles("safe-migr-transfer-ownerships-rest-predeploys");

        // Execute the full timelock flow with verification
        _rawTimelockTransaction(scheduleTx);

        // Wait for the timelock delay
        console2.log("Waiting for timelock delay");
        vm.warp(block.timestamp + MIN_DELAY + 1);

        // Execute execute transaction
        console2.log("Executing execute transaction");
        _rawTimelockTransaction(executeTx);

        // Verify that UBIPool and IPTokenStaking have new pending owner
        assertEq(ubiPool.pendingOwner(), newTimelockAddress, "UBIPool not transferred");
        assertEq(ipTokenStaking.pendingOwner(), newTimelockAddress, "IPTokenStaking not transferred");
    }

    function _testReceiveOwnershipRestPredeploys() private {
        // Run the script to generate JSON
        ReceiveOwnershipRestPredeploys script = new ReceiveOwnershipRestPredeploys();
        script.run(); // Generate operation JSON

        // Get all transaction JSONs (schedule, cancel, execute)
        (
            JSONTxWriter.Transaction memory scheduleTx,
            JSONTxWriter.Transaction memory executeTx,
            JSONTxWriter.Transaction memory cancelTx
        ) = _readAllTransactionFiles("safe-migr-receive-ownerships-rest-predeploys");

        // Execute the full timelock flow with verification
        _rawTimelockTransaction(scheduleTx);

        // Wait for the timelock delay
        console2.log("Waiting for timelock delay");
        vm.warp(block.timestamp + MIN_DELAY + 1);

        // Execute execute transaction
        console2.log("Executing execute transaction");
        _rawTimelockTransaction(executeTx);

        // Verify that the rest of the predeploys have new ownership
        assertEq(ubiPool.owner(), newTimelockAddress, "UBIPool not transferred");
        assertEq(ipTokenStaking.owner(), newTimelockAddress, "IPTokenStaking not transferred");
        assertEq(ubiPool.pendingOwner(), address(0), "UBIPool pending owner not cleared");
        assertEq(ipTokenStaking.pendingOwner(), address(0), "IPTokenStaking pending owner not cleared");
    }

    function _testRenounceOwnershipOldMultisig() private {
        // Run the script to generate JSON
        RenounceGovernanceRoles script = new RenounceGovernanceRoles();
        script.run(); // Generate operation JSON

        // Get regular transactions from JSON
        JSONTxWriter.Transaction[] memory txs = _readRegularTransactionFiles("safe-migr-renounce-gov-roles");

        // Execute the full timelock flow with verification
        for (uint256 i = 0; i < txs.length; i++) {
            console2.log("Executing transaction", i);
            _rawTimelockTransaction(txs[i]);
        }

        // Verify that the old multisig has been renounced
        assertFalse(
            newTimelock.hasRole(oldTimelock.PROPOSER_ROLE(), OLD_TIMELOCK_PROPOSER),
            "Old multisig proposer role not revoked"
        );
        assertFalse(
            newTimelock.hasRole(oldTimelock.CANCELLER_ROLE(), OLD_TIMELOCK_PROPOSER),
            "Old multisig canceller role not revoked"
        );
        assertFalse(newTimelock.hasRole(DEFAULT_ADMIN_ROLE, OLD_TIMELOCK_PROPOSER), "Old multisig admin is not admin");
        assertFalse(
            newTimelock.hasRole(oldTimelock.EXECUTOR_ROLE(), OLD_TIMELOCK_EXECUTOR),
            "Old multisig executor role not revoked"
        );
        assertFalse(
            newTimelock.hasRole(oldTimelock.CANCELLER_ROLE(), OLD_TIMELOCK_CANCELLER),
            "Old multisig canceller role not revoked"
        );

        // Verify new multisig has roles
        assertTrue(
            newTimelock.hasRole(newTimelock.DEFAULT_ADMIN_ROLE(), SAFE_TIMELOCK_PROPOSER),
            "New Safe admin is not admin"
        );
        assertTrue(
            newTimelock.hasRole(newTimelock.PROPOSER_ROLE(), SAFE_TIMELOCK_PROPOSER),
            "Old multisig proposer role not revoked"
        );
        assertTrue(
            newTimelock.hasRole(newTimelock.CANCELLER_ROLE(), SAFE_TIMELOCK_PROPOSER),
            "Old multisig canceller role not revoked"
        );
        assertTrue(
            newTimelock.hasRole(newTimelock.EXECUTOR_ROLE(), SAFE_TIMELOCK_EXECUTOR),
            "Old multisig executor role not revoked"
        );
        assertTrue(
            newTimelock.hasRole(newTimelock.CANCELLER_ROLE(), SAFE_TIMELOCK_CANCELLER),
            "Old multisig canceller role not revoked"
        );
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
}
