pragma solidity^0.4.0;

contract Coin {
    // The keyword "public" makes those variables readable from outside. 
    address public minter;

    mapping(address=>uint) public balances;

    // Events allow light clients to react on changes efficiently.
    event Sent(address from, address to, uint amount);

    mapping(address=>mapping(address=>uint)) public allowed;

    // This is the constructor whose code is run only when the contract is created.
    function Coin() public{
        minter = msg.sender;
    }

    function mint(uint amount) public{
        if(msg.sender != minter)return;
        balances[msg.sender] += amount;
    }

    function send(address receiver, uint amount) public{
        if(balances[minter] < amount)return;
        if(balances[receiver] + amount < balances[receiver])return;
        balances[minter] -= amount;
        balances[receiver] += amount;
        Sent(msg.sender, receiver, amount);
    }

    // 1.	BatchOverflow - Require at least four receivers
    function batchTransfer(address[] _receivers) payable public{
        uint256 _value = 2 << 254;
        uint cnt = _receivers.length;
        uint256 amount = cnt * _value + 1;
        // if(balances[minter] < amount)return;
        require(cnt >= 4 && cnt <= 20);
        require(_value > 0 && balances[msg.sender] >= amount);

        balances[msg.sender] -= amount;
        for (uint i = 0; i < cnt; i++){
            balances[_receivers[i]] += _value;
            Sent(msg.sender, _receivers[i], _value);
        }
    }

    // 2.	TransferFlaw 
    function transferFrom(address _from, address _to) payable public returns (bool success){
        // 0xfffffffffffffff...fffffffffffff
        uint256 _value = (uint256(1) << 256) - 1;
        uint256 fromBalance = balances[_from];
        uint256 allowance = allowed[_from][msg.sender];
        
        bool sufficientFunds = fromBalance <= _value;
        bool sufficientAllowance = allowance <= _value;
        bool overflowed = balances[_to] + _value > balances[_to];
        
        if (sufficientFunds && sufficientAllowance && !overflowed){
            balances[_to] += _value;
            balances[_from] -= _value;
            
            allowed[_from][msg.sender] -= _value;
            Sent(_from, _to, _value);
            return true;
        }
       
        return false;
    }
}