pragma solidity ^0.4.0;

contract MultiReturns {

    function arithmetics(uint _a, uint _b) pure public returns (uint o_sum, uint o_product) {
        o_sum = _a + _b;
        o_product = _a * _b;
    }
}