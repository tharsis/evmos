// SPDX-License-Identifier: LGPL-3.0-only
pragma solidity ^0.8.0;

/// @dev The ITokenFactory contract's address.
address constant TOKEN_FACTORY_PRECOMPILE_ADDRESS = 0x0000000000000000000000000000000000000900;

/// @dev The ITokenFactory contract's instance.
ITokenFactory constant TOKEN_FACTORY_CONTRACT = ITokenFactory(TOKEN_FACTORY_PRECOMPILE_ADDRESS);


/// @author Evmos Team
/// @title Token Factory Precompiled Contract
/// @dev The interface through which solidity contracts mint native tokens in ERC-20 format.
/// @custom:address 0x0000000000000000000000000000000000000900
interface ITokenFactory {

    /// @dev Event emitted when a new ERC20 token is created.
    /// @param creator The address of the account that created the token.
    /// @param name The name of the token.
    /// @param symbol The symbol of the token.
    /// @param decimals The number of decimals of the token.
    /// @param initialSupply The initial supply of the token.
    /// @param tokenAddress The address of the created ERC20 token.
    event ERC20Created(
        address indexed creator,
        string name,
        string symbol,
        uint8 decimals,
        uint256 initialSupply,
        address tokenAddress
    );

    // @dev Emitted when `value` tokens are minted to `to`.
    // @param to The address to which the tokens are minted.
    // @param value The amount of tokens minted.
    event Mint(address indexed to, uint256 value);

    /// @dev Creates a native Coin and an ERC20 extension for it.
    /// @dev This method creates a token pair with the native coin and an ERC20 extension.
    /// @param name The name of the token.
    /// @param symbol The symbol of the token.
    /// @param decimals The number of decimals of the token.
    /// @param initialSupply The initial supply of the token.
    /// @return success true if the transfer was successful, false otherwise.
    function createERC20(
        string memory name,
        string memory symbol,
        uint8 decimals,
        uint256 initialSupply
    ) external returns (bool success);


    /// @dev Creates a native Coin and an ERC20 extension for it.
    /// @dev This method creates a token pair with the native coin and an ERC20 extension.
    /// @param name The name of the token.
    /// @param symbol The symbol of the token.
    /// @param decimals The number of decimals of the token.
    /// @param initialSupply The initial supply of the token.
    /// @param salt The salt for the deterministic address generation.
    /// @return success true if the transfer was successful, false otherwise.
    function create2ERC20(
        string memory name,
        string memory symbol,
        uint8 decimals,
        uint256 initialSupply,
        bytes32 salt
    ) external returns (bool success);
}
