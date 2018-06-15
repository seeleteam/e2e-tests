pragma solidity ^0.4.0;

contract TestAddress {
    address public home;

    function set(address x) public {
        home = x;
    }
}