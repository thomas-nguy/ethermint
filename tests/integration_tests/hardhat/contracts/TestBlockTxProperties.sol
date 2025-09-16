// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

contract TestBlockTxProperties {
    function getBlockHash(uint256 blockNumber) public view returns (bytes32) {
        return blockhash(blockNumber);
    }

    function getOrigin() public view returns (address) {
        return tx.origin;
    }

    function getGasPrice() public view returns (uint256) {
        return tx.gasprice;
    }
}
