// SPDX-License-Identifier: GPL-3.0-only
pragma solidity 0.8.23;

/**
 * @title Predeploys
 * @notice Predeploy addresses (match story/genutil/evm/predeploys.go)
 */
library Predeploys {
    /// @notice Predeploys
    /// @dev Address reserved for ERC20 wrapper for Story's native token, address starts with
    /// Story mainnet chain ID
    address internal constant WIP = 0x1514000000000000000000000000000000000000;

    /// @dev We reserve the first 1024 addresses after Namespace for proxied predeploys.
    /// GenerateAlloc.s.sol will set a TransparentUpgradeableProxy for each of them, set to a
    /// deterministic implementation address. Only named predeploys's implementation will have bytecode, the
    /// others will be EOAs. Governance will be able to upgrade to valid implementations later if needed.
    address internal constant Namespace = 0xCCcCCc0000000000000000000000000000000000;
    /// @dev Number of reserved proxied predeploys in the Namespace.
    uint256 internal constant NamespaceSize = 1024;

    /// @dev IPTokenStaking proxy address
    address internal constant Staking = 0xCCcCcC0000000000000000000000000000000001;
    /// @dev UBIPool proxy address
    address internal constant UBIPool = 0xCccCCC0000000000000000000000000000000002;
    /// @dev UpgradeEntryPoint proxy address
    address internal constant Upgrades = 0xccCCcc0000000000000000000000000000000003;

    /// @notice Create3 factory address https://github.com/ZeframLou/create3-factory
    /// @dev Since Create3 is deployed using Create2, which is deterministic but depends on the deployer's wallet,
    /// we use a predeploy to maintain compatibility with the contracts deployed by ZeframLou in a permissionless way.
    address internal constant Create3 = 0x9fBB3DF7C40Da2e5A0dE984fFE2CCB7C47cd0ABf;

    /// @notice ERC6551Registry address Predeploy
    /// @dev Common address for the ERC6551Registry across all chains defined by ERC-6551
    address internal constant ERC6551Registry = 0x000000006551c19487814612e58FE06813775758;

    /// @notice Multicall3 address
    /// @dev this is not a predeploy, but it's the common address for Multicall3 across all chains
    address internal constant Multicall3 = 0xcA11bde05977b3631167028862bE2a173976CA11;

    /// @notice Return true if `addr` is proxied
    function proxied(address addr) internal pure returns (bool) {
        return addr > Namespace && addr <= address(uint160(Namespace) + uint160(NamespaceSize));
    }

    /// @notice Return implementation address for a proxied predeploy
    function getImplAddress(address proxyAddress) internal pure returns (address) {
        require(proxied(proxyAddress), "Predeploys: not proxied");

        // max uint160 is odd, which gives us unique implementation for each predeploy
        return address(type(uint160).max - uint160(proxyAddress));
    }
}
