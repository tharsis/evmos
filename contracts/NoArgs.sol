
// SPDX-License-Identifier: MIT

pragma solidity ^0.8.0;


// This is an evil token. Whenever an A -> B transfer is called, half of the amount goes to B
// and half to a predefined C
contract NoArgs {

    uint store;

    constructor(uint x) {
	store = x;
    }

    function query() external view returns(uint) {
	return store;
    }
}




