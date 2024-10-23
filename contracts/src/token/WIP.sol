// SPDX-License-Identifier: GPL-3.0-only
pragma solidity ^0.8.23;

import { ERC20 } from "solady/src/tokens/ERC20.sol";
/// @notice Wrapped IP implementation.
/// @author Inspired by WETH9 (https://github.com/dapphub/ds-weth/blob/master/src/weth9.sol)
contract WIP is ERC20 {
    /// @notice emitted when IP is deposited in exchange for WIP
    event Deposit(address indexed from, uint amount);
    /// @notice emitted when WIP is withdrawn in exchange for IP
    event Withdrawal(address indexed to, uint amount);
    /// @notice emitted when a transfer of IP fails
    error IPTransferFailed();
    /// @notice emitted when an invalid transfer recipient is detected
    error InvalidTransferReceiver();
    /// @notice emitted when an invalid transfer spender is detected
    error InvalidTransferSpender();

    /// @notice triggered when IP is deposited in exchange for WIP
    receive() external payable {
        deposit();
    }

    /// @notice deposits IP in exchange for WIP
    /// @dev the amount of IP deposited is equal to the amount of WIP minted
    function deposit() public payable {
        _mint(msg.sender, msg.value);
        emit Deposit(msg.sender, msg.value);
    }

    /// @notice withdraws WIP in exchange for IP
    /// @dev the amount of IP minted is equal to the amount of WIP burned
    /// @param value the amount of WIP to burn and withdraw
    function withdraw(uint value) external {
        _burn(msg.sender, value);
        (bool success, ) = msg.sender.call{ value: value }("");
        if (!success) {
            revert IPTransferFailed();
        }
        emit Withdrawal(msg.sender, value);
    }

    /// @notice returns the name of the token
    function name() public view override returns (string memory) {
        return "Wrapped IP";
    }

    /// @notice returns the symbol of the token
    function symbol() public view override returns (string memory) {
        return "WIP";
    }

    /// @notice approves `spender` to spend `amount` of WIP
    function approve(address spender, uint256 amount) public override returns (bool) {
        if (spender == msg.sender) {
            revert InvalidTransferSpender();
        }

        return super.approve(spender, amount);
    }

    /// @notice transfers `amount` of WIP to a recipient `to`
    function transfer(address to, uint256 amount) public override returns (bool) {
        if (to == address(0)) {
            revert InvalidTransferReceiver();
        }
        if (to == address(this)) {
            revert InvalidTransferReceiver();
        }

        return super.transfer(to, amount);
    }

    /// @notice transfers `amount` of WIP from `from` to a recipient `to`
    function transferFrom(address from, address to, uint256 amount) public override returns (bool) {
        if (to == address(0)) {
            revert InvalidTransferReceiver();
        }
        if (to == address(this)) {
            revert InvalidTransferReceiver();
        }

        return super.transferFrom(from, to, amount);
    }

    /// @dev Sets Permit2 contract's allowance to infinity.
    function _givePermit2InfiniteAllowance() internal pure override returns (bool) {
        return true;
    }
}
