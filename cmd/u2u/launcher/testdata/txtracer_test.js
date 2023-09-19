var a1 = personal.newAccount("");
personal.unlockAccount(a1, "", 300);
admin.sleep(1);
u2u.sendTransaction({ from: u2u.accounts[0], to: a1, value: "10000000000000000000" });
admin.sleep(2);
var abi = [{ "inputs": [], "name": "deploy", "outputs": [{ "internalType": "address", "name": "", "type": "address" }], "stateMutability": "nonpayable", "type": "function" }, { "inputs": [], "name": "getA", "outputs": [{ "internalType": "uint256", "name": "", "type": "uint256" }], "stateMutability": "view", "type": "function" }, { "inputs": [{ "internalType": "uint256", "name": "_a", "type": "uint256" }], "name": "setA", "outputs": [], "stateMutability": "nonpayable", "type": "function" }, { "inputs": [{ "internalType": "uint256", "name": "_a", "type": "uint256" }], "name": "setInA", "outputs": [], "stateMutability": "nonpayable", "type": "function" }, { "inputs": [], "name": "tst", "outputs": [{ "internalType": "address", "name": "", "type": "address" }], "stateMutability": "view", "type": "function" }];
var bytecode = "0x608060405234801561001057600080fd5b50610589806100206000396000f3fe608060405234801561001057600080fd5b50600436106100575760003560e01c8063775c300c1461005c57806391888f2e146100a6578063d46300fd146100f0578063ea87f3731461010e578063ee919d501461013c575b600080fd5b61006461016a565b604051808273ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200191505060405180910390f35b6100ae61032a565b604051808273ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200191505060405180910390f35b6100f8610350565b6040518082815260200191505060405180910390f35b61013a6004803603602081101561012457600080fd5b8101908080359060200190929190505050610359565b005b6101686004803603602081101561015257600080fd5b81019080803590602001909291905050506103ef565b005b6000807f746573740000000000000000000000000000000000000000000000000000000060001b607b60405161019f906103f9565b808281526020019150508190604051809103906000f59050801580156101c9573d6000803e3d6000fd5b5090508073ffffffffffffffffffffffffffffffffffffffff1663d46300fd6040518163ffffffff1660e01b815260040160206040518083038186803b15801561021257600080fd5b505afa158015610226573d6000803e3d6000fd5b505050506040513d602081101561023c57600080fd5b81019080805190602001909291905050506000819055508073ffffffffffffffffffffffffffffffffffffffff1663ee919d506101416040518263ffffffff1660e01b815260040180828152602001915050600060405180830381600087803b1580156102a857600080fd5b505af11580156102bc573d6000803e3d6000fd5b5050505080600160006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff160217905550600160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1691505090565b600160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1681565b60008054905090565b6000600160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1690508073ffffffffffffffffffffffffffffffffffffffff1663ee919d50836040518263ffffffff1660e01b815260040180828152602001915050600060405180830381600087803b1580156103d357600080fd5b505af11580156103e7573d6000803e3d6000fd5b505050505050565b8060008190555050565b61014d806104078339019056fe608060405234801561001057600080fd5b5060405161014d38038061014d8339818101604052602081101561003357600080fd5b8101908080519060200190929190505050806000819055505060f38061005a6000396000f3fe6080604052348015600f57600080fd5b5060043610603c5760003560e01c80630dbe671f146041578063d46300fd14605d578063ee919d50146079575b600080fd5b604760a4565b6040518082815260200191505060405180910390f35b606360aa565b6040518082815260200191505060405180910390f35b60a260048036036020811015608d57600080fd5b810190808035906020019092919050505060b3565b005b60005481565b60008054905090565b806000819055505056fea264697066735822122084a43fd050d3de9bc0cdc2f86b02db81bc69b5639cedd1a5b72010d2a664879a64736f6c63430006020033a2646970667358221220a00266288222dc5d9803c3840f0e4d5d04630cb6bac841d30c42511d9b58045864736f6c63430006020033";
var simpleContract = u2u.contract(abi);
var simpleContractAddr = ""
var simple = simpleContract.new(42, { from: u2u.accounts[0], data: bytecode, gas: 0x47b760 }, function (e, contract) {
    if (!e) {
        if (contract.address) {
            simpleContractAddr = contract.address
            // Initialize test contract for interaction
            var testContract = u2u.contract(abi)
            testContract = testContract.at(simpleContractAddr)

            // Call simple contract call to check created trace
            var tx = testContract.setA.sendTransaction(24, { from: u2u.accounts[0] })
            var txDeploy = testContract.deploy.sendTransaction({ from: u2u.accounts[0] })
            admin.sleep(1);

            console.log(tx)
            console.log(txDeploy)
        }
    }
});
admin.sleep(3);

/*
// Contract source code:
pragma solidity 0.6.2;

contract DeployTest {
    uint256 a;
    address public tst;
    
    function deploy() public returns(address) {
        // Creating new contract using CREATE2 instruction
        Test test = new Test{salt: 0x7465737400000000000000000000000000000000000000000000000000000000}(123);
        a = test.getA();
        test.setA(321);
        tst = address(test);
        return tst;
    }
    
    function setA(uint256 _a) public {
        a = _a;
    }
    
    function getA() public view returns (uint256) {
        return a;
    }
    
    function setInA(uint256 _a) public {
        Test t = Test(tst);
        t.setA(_a);
    }
}

contract Test {
    uint256 public a;
    constructor (uint256 _a) public {
        a = _a;
    }
    
    function setA(uint256 _a) public {
        a = _a;
    }
    
    function getA() public view returns (uint256) {
        return a;
    }
}
*/