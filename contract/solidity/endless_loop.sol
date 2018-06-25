pragma solidity ^0.4.0;

contract C {
    function g(uint) public returns (uint ret) { return f(); }
    function f() public returns (uint ret) { return g(7) + f(); }
    
    function l() pure public {
        while (true) {
        } 
    }
}