// SPDX-License-Identifier: LGPL-3.0-only
pragma solidity >=0.8.18;

/**
 * @author Evmos Team
 * @title Wrapped ERC20 Interface
 * @dev Interface for representing the native EVM token as ERC20 standard.
 */
interface IWERC20 {
    /// @dev Emitted when the native tokens are deposited in exchange for the wrapped ERC20.
    /// @param dst The account for which the deposit is made.
    /// @param wad The amount of native tokens deposited.
    event Deposit(address indexed dst, uint256 wad);

    /// @dev Emitted when the native token is withdrawn.
    /// @param src The account for which the withdrawal is made.
    /// @param wad The amount of native tokens withdrawn.
    event Withdrawal(address indexed src, uint256 wad);

    /// @dev Default fallback payable function. Must call the deposit method in implementing contracts.
    fallback() external payable;

    /// @dev Default receive payable function. Must call the deposit method in implementing contracts.
    receive() external payable;

    /// @dev Deposits native tokens in exchange for wrapped ERC20 token.
    /// @dev After execution of this function the SetBalance function
    /// @dev burns the tokens and increases the contract balance of the ERC20 tokens.
    /// @dev Emits a Deposit Event.
    function deposit() external payable;

    /// @dev Withdraws native tokens from wrapped ERC20 token.
    /// @dev After execution of this function the SetBalance function
    /// @dev decreases the contract balance of the ERC20 tokens and mints the tokens.
    /// @dev Emits a Withdrawal Event.
    /// @param wad The amount of native tokens to be withdrawn.
    function withdraw(uint256 wad) external;
}