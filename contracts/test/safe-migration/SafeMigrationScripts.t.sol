// SPDX-License-Identifier: MIT
pragma solidity 0.8.23;

import { console2 } from "forge-std/Test.sol";
import { Test} from "test/utils/Test.sol";
import { Vm } from "forge-std/Vm.sol";
import { stdJson } from "forge-std/StdJson.sol";

import { TimelockController } from "@openzeppelin/contracts/governance/TimelockController.sol";
import { Ownable } from "@openzeppelin/contracts/access/Ownable.sol";
import { ProxyAdmin } from "@openzeppelin/contracts/proxy/transparent/ProxyAdmin.sol";
import { Ownable2StepUpgradeable } from "@openzeppelin/contracts-upgradeable/access/Ownable2StepUpgradeable.sol";
import { AccessControl } from "@openzeppelin/contracts/access/AccessControl.sol";

import { Predeploys } from "src/libraries/Predeploys.sol";
import { Create3 } from "src/deploy/Create3.sol";
import { JSONTxWriter } from "script/utils/JSONTxWriter.s.sol";

// Import migration scripts
import { DeployNewTimelock } from "script/admin-actions/migrate-to-safe/1.DeployNewTimelock.s.sol";
import { TransferOwnershipsProxyAdmin1 } from "script/admin-actions/migrate-to-safe/2.1.TransferOwnershipProxyAdmin1.s.sol";
import { TransferOwnershipsProxyAdmin2 } from "script/admin-actions/migrate-to-safe/2.2.TransferOwnershipProxyAdmin2.s.sol";
// import { TransferOwnershipsUpgradesEntrypoint } from "script/admin-actions/migrate-to-safe/3.1.TransferOwnershipUpgradesEntrypoint.s.sol";
// import { ReceiveOwnershipUpgradesEntrypoint } from "script/admin-actions/migrate-to-safe/3.2.ReceiveOwnershipUpgradesEntryPoint.s.sol";
// import { TransferOwnershipsRestPredeploys } from "script/admin-actions/migrate-to-safe/3.3.TransferOwnershipRestPredeploys.s.sol";
// import { ReceiveOwnershipRestPredeploys } from "script/admin-actions/migrate-to-safe/3.4.ReceiveOwnershipRestPredeploys.s.sol";
// import { RenounceOldMultisigRoles } from "script/admin-actions/migrate-to-safe/4.RenounceOwnershipOldMultisig.sol";

contract SafeMigrationScriptsTest is Test {
    using stdJson for string;

    // Mock addresses for the test
    address private OLD_TIMELOCK_PROPOSER;
    address private OLD_TIMELOCK_GUARDIAN;
    address private OLD_EXECUTOR_ADDRESS;
    address private constant SAFE_TIMELOCK_PROPOSER = address(0x1111111111111111111111111111111111111111);
    address private constant SAFE_TIMELOCK_EXECUTOR_ADDRESS = address(0x2222222222222222222222222222222222222222);
    address private constant SAFE_TIMELOCK_GUARDIAN_ADDRESS = address(0x3333333333333333333333333333333333333333);

    // Private key for timelock deployer
    uint256 private constant DEPLOYER_PRIVATE_KEY = 0x1;

    // Timelock controller addresses
    address private oldTimelockAddress;
    address private newTimelockAddress;
    TimelockController private oldTimelock;
    TimelockController private newTimelock;

    // Mock contracts for testing
    address[] private proxyAdmins;
    // JSON directories
    string constant TIMELOCK_ROOT_DIR = "./broadcast/";
    string constant OUTPUT_DIR = "./output/";

    // Role constants from TimelockController
    bytes32 private constant DEFAULT_ADMIN_ROLE = 0x00;

    function setUp() public override {
        super.setUp();
        OLD_TIMELOCK_PROPOSER = admin;
        OLD_TIMELOCK_GUARDIAN = guardian;
        
        oldTimelock = timelock;

        // Setup environment variables
        _setupEnvVars();
    }
    
    function testMigrationScripts() public {
        // Step 1: Deploy new timelock
        _runDeployNewTimelock();
        /*
        // Step 2.1: Transfer ownership of first half of proxy admins
        _runTransferOwnershipProxyAdmin1();
        
        // Step 2.2: Transfer ownership of second half of proxy admins
        _runTransferOwnershipProxyAdmin2();
        
        // Step 3.1: Transfer ownership of UpgradesEntrypoint
        _runTransferOwnershipUpgradesEntrypoint();
        
        // Step 3.2: Accept ownership of UpgradesEntrypoint
        _runAcceptOwnershipUpgradesEntrypoint();
        
        // Step 3.3: Transfer ownership of rest of predeploys
        _runTransferOwnershipRestPredeploys();
        
        // Step 3.4: Accept ownership of rest of predeploys
        _runAcceptOwnershipRestPredeploys();
        
        // Step 4: Renounce ownership of old multisig
        _runRenounceOwnershipOldMultisig();
        
        // Final verification that all permissions are correctly set
        _verifyFinalState();
        */
    }
    
    function _setupEnvVars() private {
        // Setup environment variables
        vm.setEnv("NEW_TIMELOCK_DEPLOYER_PRIVATE_KEY", vm.toString(DEPLOYER_PRIVATE_KEY));
        vm.setEnv("SAFE_TIMELOCK_PROPOSER", vm.toString(SAFE_TIMELOCK_PROPOSER));
        vm.label(SAFE_TIMELOCK_PROPOSER, "SafeTimelockProposer");
        vm.setEnv("SAFE_TIMELOCK_EXECUTOR_ADDRESS", vm.toString(SAFE_TIMELOCK_EXECUTOR_ADDRESS));
        vm.label(SAFE_TIMELOCK_EXECUTOR_ADDRESS, "SafeTimelockExecutor");
        vm.setEnv("SAFE_TIMELOCK_GUARDIAN_ADDRESS", vm.toString(SAFE_TIMELOCK_GUARDIAN_ADDRESS));
        vm.label(SAFE_TIMELOCK_GUARDIAN_ADDRESS, "SafeTimelockGuardian");
        vm.setEnv("OLD_TIMELOCK_ADDRESS", vm.toString(address(oldTimelock)));
        vm.label(address(oldTimelock), "OldTimelock");
        vm.setEnv("OLD_TIMELOCK_PROPOSER", vm.toString(OLD_TIMELOCK_PROPOSER));
        vm.label(OLD_TIMELOCK_PROPOSER, "Old Timelock Proposer");
        vm.setEnv("OLD_TIMELOCK_GUARDIAN_ADDRESS", vm.toString(OLD_TIMELOCK_GUARDIAN));
        vm.label(OLD_TIMELOCK_GUARDIAN, "Old Timelock Guardian");
        vm.setEnv("MIN_DELAY", "300"); // 5 minutes
    }
    
    function _runDeployNewTimelock() private {
        DeployNewTimelock deployScript = new DeployNewTimelock();
        deployScript.run();
        
        // Get new timelock address
        bytes32 salt = keccak256("STORY_TIMELOCK_CONTROLLER_SAFE");
        newTimelockAddress = Create3(Predeploys.Create3).predictDeterministicAddress(salt);
        newTimelock = TimelockController(payable(newTimelockAddress));
        
        // Verify new timelock setup
        assertTrue(newTimelockAddress.code.length > 0, "New timelock not deployed");
        assertTrue(newTimelock.hasRole(newTimelock.PROPOSER_ROLE(), SAFE_TIMELOCK_PROPOSER), "Safe admin not proposer");
        assertTrue(newTimelock.hasRole(newTimelock.PROPOSER_ROLE(), OLD_TIMELOCK_PROPOSER), "Old multisig not proposer");
        assertTrue(newTimelock.hasRole(newTimelock.EXECUTOR_ROLE(), SAFE_TIMELOCK_EXECUTOR_ADDRESS), "Safe executor not executor");
        assertTrue(newTimelock.hasRole(newTimelock.EXECUTOR_ROLE(), OLD_TIMELOCK_PROPOSER), "Old multisig not executor");
        assertTrue(newTimelock.hasRole(newTimelock.CANCELLER_ROLE(), SAFE_TIMELOCK_GUARDIAN_ADDRESS), "Safe guardian not canceller");
        assertTrue(newTimelock.hasRole(newTimelock.CANCELLER_ROLE(), OLD_TIMELOCK_GUARDIAN), "Old security council not canceller");
        assertTrue(newTimelock.hasRole(DEFAULT_ADMIN_ROLE, SAFE_TIMELOCK_PROPOSER), "New Safe admin is not admin");
        assertTrue(newTimelock.hasRole(DEFAULT_ADMIN_ROLE, OLD_TIMELOCK_PROPOSER), "Old multisig admin is not admin");
        assertFalse(newTimelock.hasRole(DEFAULT_ADMIN_ROLE, deployer), "deployer is admin");
        assertEq(newTimelock.getMinDelay(), 300, "Delay not set correctly");
    }
    
    /*
    function _runTransferOwnershipProxyAdmin1() private {
        // Run the script to generate JSON
        TransferOwnershipsProxyAdmin1 script = new TransferOwnershipsProxyAdmin1();
        script.generate(); // Generate operation JSON
        
        // Get the JSON output
        string memory jsonPath = string.concat(OUTPUT_DIR, "safe-migration-transfer-ownerships-proxy-admin-2", ".json");
        string memory json = vm.readFile(jsonPath);
        
        // Execute transactions for old timelock to transfer ownership
        _executeOldTimelockBatchTransaction(json);
        
        // Verify that half of the ProxyAdmins have new ownership
        uint256 halfSize = Predeploys.NamespaceSize / 2;
        uint256 startIdx = halfSize; // Second half of proxies
        
        for (uint160 i = uint160(startIdx); i < Predeploys.NamespaceSize; i++) {
            ProxyAdmin proxyAdmin = ProxyAdmin(proxyAdmins[i]);
            assertEq(proxyAdmin.owner(), newTimelockAddress, "ProxyAdmin ownership not transferred");
        }
    }
    
    function _runTransferOwnershipProxyAdmin2() private {
        // Run the script to generate JSON
        TransferOwnershipsProxyAdmin2 script = new TransferOwnershipsProxyAdmin2();
        script.generate(); // Generate operation JSON
        
        // Get the JSON output
        string memory jsonPath = string.concat(OUTPUT_DIR, "safe-migration-transfer-ownerships-proxy-admin-1", ".json");
        string memory json = vm.readFile(jsonPath);
        
        // Execute transactions for old timelock to transfer ownership
        _executeOldTimelockBatchTransaction(json);
        
        // Verify that all ProxyAdmins have new ownership now
        for (uint160 i = 0; i < Predeploys.NamespaceSize; i++) {
            ProxyAdmin proxyAdmin = ProxyAdmin(proxyAdmins[i]);
            assertEq(proxyAdmin.owner(), newTimelockAddress, "ProxyAdmin ownership not transferred");
        }
    }
    
    function _runTransferOwnershipUpgradesEntrypoint() private {
        // Run the script to generate JSON
        TransferOwnershipsUpgradesEntrypoint script = new TransferOwnershipsUpgradesEntrypoint();
        script.generate(); // Generate operation JSON
        
        // Get the JSON output
        string memory jsonPath = string.concat(OUTPUT_DIR, "safe-migration-transfer-ownerships-upgradeEntrypoint-entrypoint", ".json");
        string memory json = vm.readFile(jsonPath);
        
        // Execute transactions for old timelock to transfer ownership
        _executeOldTimelockTransaction(json);
        
        // Verify pending ownership
        assertEq(MockOwnable2StepUpgradeable(Predeploys.Upgrades).pendingOwner(), newTimelockAddress, "Pending owner not set");
        assertEq(MockOwnable2StepUpgradeable(Predeploys.Upgrades).owner(), oldTimelockAddress, "Owner should still be old timelock");
    }
    
    function _runAcceptOwnershipUpgradesEntrypoint() private {
        // Run the script to generate JSON
        ReceiveOwnershipUpgradesEntrypoint script = new ReceiveOwnershipUpgradesEntrypoint();
        script.generate(); // Generate operation JSON
        
        // Get the JSON output
        string memory jsonPath = string.concat(OUTPUT_DIR, "safe-migration-receive-ownerships-upgradeEntrypoint-entrypoint", ".json");
        string memory json = vm.readFile(jsonPath);
        
        // Execute transactions for new timelock to accept ownership
        _executeNewTimelockTransaction(json);
        
        // Verify ownership transfer
        assertEq(MockOwnable2StepUpgradeable(Predeploys.Upgrades).owner(), newTimelockAddress, "Ownership not transferred to new timelock");
        
        // Verify old timelock can't call restricted functions
        vm.startPrank(address(oldTimelock));
        vm.expectRevert("Ownable: caller is not the owner");
        MockOwnable2StepUpgradeable(Predeploys.Upgrades).restrictedFunction();
        vm.stopPrank();
        
        // Verify new timelock can call restricted functions
        vm.startPrank(address(newTimelock));
        MockOwnable2StepUpgradeable(Predeploys.Upgrades).restrictedFunction();
        vm.stopPrank();
    }
    
    function _runTransferOwnershipRestPredeploys() private {
        // Run the script to generate JSON
        TransferOwnershipRestPredeploys script = new TransferOwnershipRestPredeploys();
        script.generate(); // Generate operation JSON
        
        // Get the JSON output
        string memory jsonPath = string.concat(OUTPUT_DIR, "safe-migration-transfer-ownerships-rest-predeploys", ".json");
        string memory json = vm.readFile(jsonPath);
        
        // Execute transactions for old timelock to transfer ownership
        _executeOldTimelockBatchTransaction(json);
        
        // Verify pending ownership
        assertEq(MockOwnable2StepUpgradeable(Predeploys.Staking).pendingOwner(), newTimelockAddress, "Pending owner not set for Staking");
        assertEq(MockOwnable2StepUpgradeable(Predeploys.UBIPool).pendingOwner(), newTimelockAddress, "Pending owner not set for UBIPool");
    }
    
    function _runAcceptOwnershipRestPredeploys() private {
        // Run the script to generate JSON
        ReceiveOwnershipRestPredeploys script = new ReceiveOwnershipRestPredeploys();
        script.generate(); // Generate operation JSON
        
        // Get the JSON output
        string memory jsonPath = string.concat(OUTPUT_DIR, "safe-migration-receive-ownerships-rest-predeploys", ".json");
        string memory json = vm.readFile(jsonPath);
        
        // Execute transactions for new timelock to accept ownership
        _executeNewTimelockBatchTransaction(json);
        
        // Verify ownership transfer
        assertEq(MockOwnable2StepUpgradeable(Predeploys.Staking).owner(), newTimelockAddress, "Ownership not transferred for Staking");
        assertEq(MockOwnable2StepUpgradeable(Predeploys.UBIPool).owner(), newTimelockAddress, "Ownership not transferred for UBIPool");
    }
    
    function _runRenounceOwnershipOldMultisig() private {
        // Run the script to generate JSON
        RenounceOldMultisigRoles script = new RenounceOldMultisigRoles();
        script.run(); // Generate operation JSON
        
        // Get the JSON output
        string memory basePath = string.concat("./script/admin-actions/output/", vm.toString(block.chainid), "/");
        string memory jsonPath = string.concat(basePath, "renounce-old-multisig-roles-execute.json");
        string memory json = vm.readFile(jsonPath);
        
        // Parse the JSON array
        uint256 txCount = stdJson.readUint(json, ".length");
        console2.log("Found", txCount, "transactions in JSON");
        
        // Execute each transaction in the JSON
        for (uint256 i = 0; i < txCount; i++) {
            string memory txIndexPath = string.concat("[", vm.toString(i), "]");
            address from = stdJson.readAddress(json, string.concat(txIndexPath, ".from"));
            address to = stdJson.readAddress(json, string.concat(txIndexPath, ".to"));
            bytes memory data = stdJson.readBytes(json, string.concat(txIndexPath, ".data"));
            string memory comment = stdJson.readString(json, string.concat(txIndexPath, ".comment"));
            
            console2.log("Executing transaction:", comment);
            console2.log("From:", from);
            console2.log("To:", to);
            
            // Execute the transaction
            vm.startPrank(from);
            (bool success, ) = to.call(data);
            require(success, "Transaction execution failed");
            vm.stopPrank();
        }
        
        // Verify roles revoked
        assertFalse(newTimelock.hasRole(newTimelock.PROPOSER_ROLE(), OLD_TIMELOCK_PROPOSER), "Old multisig still has proposer role");
        assertFalse(newTimelock.hasRole(newTimelock.EXECUTOR_ROLE(), OLD_TIMELOCK_PROPOSER), "Old multisig still has executor role");
        assertFalse(newTimelock.hasRole(newTimelock.CANCELLER_ROLE(), OLD_TIMELOCK_GUARDIAN), "Old security council still has canceller role");
        assertFalse(newTimelock.hasRole(DEFAULT_ADMIN_ROLE, OLD_TIMELOCK_PROPOSER), "Old multisig still has admin role");
    }
    
    function _verifyFinalState() private {
        // Verify all ProxyAdmins are owned by new timelock
        for (uint160 i = 0; i < proxyAdmins.length; i++) {
            ProxyAdmin proxyAdmin = ProxyAdmin(proxyAdmins[i]);
            assertEq(proxyAdmin.owner(), newTimelockAddress, "ProxyAdmin not owned by new timelock");
        }
        
        // Verify all predeploys are owned by new timelock
        assertEq(MockOwnable2StepUpgradeable(Predeploys.Upgrades).owner(), newTimelockAddress, "Upgrades not owned by new timelock");
        assertEq(MockOwnable2StepUpgradeable(Predeploys.Staking).owner(), newTimelockAddress, "Staking not owned by new timelock");
        assertEq(MockOwnable2StepUpgradeable(Predeploys.UBIPool).owner(), newTimelockAddress, "UBIPool not owned by new timelock");
        
        // Verify all old owners can't call restricted functions
        vm.startPrank(address(oldTimelock));
        vm.expectRevert("Ownable: caller is not the owner");
        MockOwnable2StepUpgradeable(Predeploys.Upgrades).restrictedFunction();
        vm.expectRevert("Ownable: caller is not the owner");
        MockOwnable2StepUpgradeable(Predeploys.Staking).restrictedFunction();
        vm.expectRevert("Ownable: caller is not the owner");
        MockOwnable2StepUpgradeable(Predeploys.UBIPool).restrictedFunction();
        vm.stopPrank();
        
        // Verify new owner can call restricted functions
        vm.startPrank(address(newTimelock));
        MockOwnable2StepUpgradeable(Predeploys.Upgrades).restrictedFunction();
        MockOwnable2StepUpgradeable(Predeploys.Staking).restrictedFunction();
        MockOwnable2StepUpgradeable(Predeploys.UBIPool).restrictedFunction();
        vm.stopPrank();
    }
    
    // Helper functions to execute timelock operations
    function _executeOldTimelockTransaction(string memory json) private {
        // Parse the Transaction from JSON
        address target = stdJson.readAddress(json, ".target");
        uint256 value = stdJson.readUint(json, ".value");
        bytes memory data = stdJson.readBytes(json, ".data");
        bytes32 predecessor = stdJson.readBytes32(json, ".predecessor");
        bytes32 salt = stdJson.readBytes32(json, ".salt");
        uint256 delay = stdJson.readUint(json, ".delay");
        
        // Schedule the transaction in the old timelock
        vm.startPrank(OLD_TIMELOCK_PROPOSER);
        bytes32 operationId = oldTimelock.hashOperation(
            target,
            value,
            data,
            predecessor,
            salt
        );
        
        oldTimelock.schedule(
            target,
            value,
            data,
            predecessor,
            salt,
            delay
        );
        
        // Fast forward past the delay
        vm.warp(block.timestamp + delay + 1);
        
        // Execute the transaction
        oldTimelock.execute(
            target,
            value,
            data,
            predecessor,
            salt
        );
        vm.stopPrank();
    }
    
    function _executeOldTimelockBatchTransaction(string memory json) private {
        // Parse the array of targets, values, and data
        uint256 length = stdJson.readUint(json, ".targets.length");
        address[] memory targets = new address[](length);
        uint256[] memory values = new uint256[](length);
        bytes[] memory datas = new bytes[](length);
        
        for (uint256 i = 0; i < length; i++) {
            string memory indexStr = vm.toString(i);
            targets[i] = stdJson.readAddress(json, string.concat(".targets[", indexStr, "]"));
            values[i] = stdJson.readUint(json, string.concat(".values[", indexStr, "]"));
            datas[i] = stdJson.readBytes(json, string.concat(".datas[", indexStr, "]"));
        }
        
        bytes32 predecessor = stdJson.readBytes32(json, ".predecessor");
        bytes32 salt = stdJson.readBytes32(json, ".salt");
        uint256 delay = stdJson.readUint(json, ".delay");
        
        // Schedule the batch transaction in the old timelock
        vm.startPrank(OLD_TIMELOCK_PROPOSER);
        bytes32 operationId = oldTimelock.hashOperationBatch(
            targets,
            values,
            datas,
            predecessor,
            salt
        );
        
        oldTimelock.scheduleBatch(
            targets,
            values,
            datas,
            predecessor,
            salt,
            delay
        );
        
        // Fast forward past the delay
        vm.warp(block.timestamp + delay + 1);
        
        // Execute the transaction
        oldTimelock.executeBatch(
            targets,
            values,
            datas,
            predecessor,
            salt
        );
        vm.stopPrank();
    }
    
    function _executeNewTimelockTransaction(string memory json) private {
        // Parse the Transaction from JSON
        address target = stdJson.readAddress(json, ".target");
        uint256 value = stdJson.readUint(json, ".value");
        bytes memory data = stdJson.readBytes(json, ".data");
        bytes32 predecessor = stdJson.readBytes32(json, ".predecessor");
        bytes32 salt = stdJson.readBytes32(json, ".salt");
        uint256 delay = stdJson.readUint(json, ".delay");
        
        // Schedule the transaction in the new timelock
        vm.startPrank(SAFE_TIMELOCK_PROPOSER);
        bytes32 operationId = newTimelock.hashOperation(
            target,
            value,
            data,
            predecessor,
            salt
        );
        
        newTimelock.schedule(
            target,
            value,
            data,
            predecessor,
            salt,
            delay
        );
        
        // Fast forward past the delay
        vm.warp(block.timestamp + delay + 1);
        
        // Execute the transaction
        vm.startPrank(SAFE_TIMELOCK_EXECUTOR_ADDRESS);
        newTimelock.execute(
            target,
            value,
            data,
            predecessor,
            salt
        );
        vm.stopPrank();
    }
    
    function _executeNewTimelockBatchTransaction(string memory json) private {
        // Parse the array of targets, values, and data
        uint256 length = stdJson.readUint(json, ".targets.length");
        address[] memory targets = new address[](length);
        uint256[] memory values = new uint256[](length);
        bytes[] memory datas = new bytes[](length);
        
        for (uint256 i = 0; i < length; i++) {
            string memory indexStr = vm.toString(i);
            targets[i] = stdJson.readAddress(json, string.concat(".targets[", indexStr, "]"));
            values[i] = stdJson.readUint(json, string.concat(".values[", indexStr, "]"));
            datas[i] = stdJson.readBytes(json, string.concat(".datas[", indexStr, "]"));
        }
        
        bytes32 predecessor = stdJson.readBytes32(json, ".predecessor");
        bytes32 salt = stdJson.readBytes32(json, ".salt");
        uint256 delay = stdJson.readUint(json, ".delay");
        
        // Schedule the batch transaction in the new timelock
        vm.startPrank(SAFE_TIMELOCK_PROPOSER);
        bytes32 operationId = newTimelock.hashOperationBatch(
            targets,
            values,
            datas,
            predecessor,
            salt
        );
        
        newTimelock.scheduleBatch(
            targets,
            values,
            datas,
            predecessor,
            salt,
            delay
        );
        
        // Fast forward past the delay
        vm.warp(block.timestamp + delay + 1);
        
        // Execute the transaction
        vm.startPrank(SAFE_TIMELOCK_EXECUTOR_ADDRESS);
        newTimelock.executeBatch(
            targets,
            values,
            datas,
            predecessor,
            salt
        );
        vm.stopPrank();
    }
    */
}

// Mock contract for testing
contract MockOwnable2StepUpgradeable is Ownable2StepUpgradeable {
    bool public restrictedFunctionCalled;
    
    function initialize(address initialOwner) public initializer {
        __Ownable_init(initialOwner);
    }
    
    function restrictedFunction() public onlyOwner {
        restrictedFunctionCalled = true;
    }
}
