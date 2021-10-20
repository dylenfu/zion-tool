/*
 * Copyright (C) 2021 The Zion Authors
 * This file is part of The Zion library.
 *
 * The Zion is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Lesser General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * The Zion is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU Lesser General Public License for more details.
 *
 * You should have received a copy of the GNU Lesser General Public License
 * along with The Zion.  If not, see <http://www.gnu.org/licenses/>.
 */

// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.
package neo_proof

import (
	"math/big"
	"strings"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/event"
)

// Reference imports to suppress errors if they are not otherwise used.
var (
	_ = big.NewInt
	_ = strings.NewReader
	_ = ethereum.NotFound
	_ = bind.Bind
	_ = common.Big1
	_ = types.BloomLookup
	_ = event.NewSubscription
)

var (
	MethodRemove = "remove"

	MethodSet = "set"

	MethodSetGov = "setGov"

	MethodGet = "get"

	MethodGovAddress = "govAddress"

	MethodProofName = "proofName"

	MethodProofs = "proofs"

	EventRemove = "Remove"

	EventSet = "Set"

	EventSetGov = "SetGov"
)

// ProofABI is the input ABI used to generate the binding from.
const ProofABI = "[{\"inputs\":[{\"internalType\":\"string\",\"name\":\"_proofName\",\"type\":\"string\"},{\"internalType\":\"address\",\"name\":\"_govAddress\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"string\",\"name\":\"_key\",\"type\":\"string\"}],\"name\":\"Remove\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"string\",\"name\":\"_key\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"_value\",\"type\":\"string\"}],\"name\":\"Set\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"_govAddress\",\"type\":\"address\"}],\"name\":\"SetGov\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"_key\",\"type\":\"string\"}],\"name\":\"get\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"govAddress\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"proofName\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"name\":\"proofs\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"_key\",\"type\":\"string\"}],\"name\":\"remove\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"_key\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"_value\",\"type\":\"string\"}],\"name\":\"set\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_govAddress\",\"type\":\"address\"}],\"name\":\"setGov\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]"

// ProofFuncSigs maps the 4-byte function signature to its string representation.
var ProofFuncSigs = map[string]string{
	"693ec85e": "get(string)",
	"46008a07": "govAddress()",
	"27d2fce5": "proofName()",
	"3b18f691": "proofs(string)",
	"80599e4b": "remove(string)",
	"e942b516": "set(string,string)",
	"cfad57a2": "setGov(address)",
}

// ProofBin is the compiled bytecode used for deploying new contracts.
var ProofBin = "0x60806040523480156200001157600080fd5b5060405162000a5138038062000a5183398101604081905262000034916200014a565b81516200004990600090602085019062000071565b50600180546001600160a01b0319166001600160a01b03929092169190911790555062000278565b8280546200007f906200023b565b90600052602060002090601f016020900481019282620000a35760008555620000ee565b82601f10620000be57805160ff1916838001178555620000ee565b82800160010185558215620000ee579182015b82811115620000ee578251825591602001919060010190620000d1565b50620000fc92915062000100565b5090565b5b80821115620000fc576000815560010162000101565b634e487b7160e01b600052604160045260246000fd5b80516001600160a01b03811681146200014557600080fd5b919050565b600080604083850312156200015e57600080fd5b82516001600160401b03808211156200017657600080fd5b818501915085601f8301126200018b57600080fd5b815181811115620001a057620001a062000117565b604051601f8201601f19908116603f01168101908382118183101715620001cb57620001cb62000117565b81604052828152602093508884848701011115620001e857600080fd5b600091505b828210156200020c5784820184015181830185015290830190620001ed565b828211156200021e5760008484830101525b9550620002309150508582016200012d565b925050509250929050565b600181811c908216806200025057607f821691505b602082108114156200027257634e487b7160e01b600052602260045260246000fd5b50919050565b6107c980620002886000396000f3fe608060405234801561001057600080fd5b506004361061007d5760003560e01c8063693ec85e1161005b578063693ec85e146100de57806380599e4b146100f1578063cfad57a214610106578063e942b5161461011957600080fd5b806327d2fce5146100825780633b18f691146100a057806346008a07146100b3575b600080fd5b61008a61012c565b6040516100979190610569565b60405180910390f35b61008a6100ae366004610626565b6101ba565b6001546100c6906001600160a01b031681565b6040516001600160a01b039091168152602001610097565b61008a6100ec366004610626565b6101de565b6101046100ff366004610626565b61028e565b005b610104610114366004610663565b610325565b61010461012736600461068c565b61039d565b60008054610139906106f0565b80601f0160208091040260200160405190810160405280929190818152602001828054610165906106f0565b80156101b25780601f10610187576101008083540402835291602001916101b2565b820191906000526020600020905b81548152906001019060200180831161019557829003601f168201915b505050505081565b805160208183018101805160028252928201919093012091528054610139906106f0565b60606002826040516101f0919061072b565b90815260200160405180910390208054610209906106f0565b80601f0160208091040260200160405190810160405280929190818152602001828054610235906106f0565b80156102825780601f1061025757610100808354040283529160200191610282565b820191906000526020600020905b81548152906001019060200180831161026557829003601f168201915b50505050509050919050565b6001546001600160a01b031633146102c15760405162461bcd60e51b81526004016102b890610747565b60405180910390fd5b6002816040516102d1919061072b565b908152602001604051809103902060006102eb9190610437565b7f834a2d47e948021d7136fb7275b3f1e1feae6333c0d683e8c13f901667defd8c8160405161031a9190610569565b60405180910390a150565b6001546001600160a01b0316331461034f5760405162461bcd60e51b81526004016102b890610747565b600180546001600160a01b0319166001600160a01b0383169081179091556040519081527f91a8c1cc2d4a3bb60738481947a00cbb9899c822916694cf8bb1d68172fdcd549060200161031a565b6001546001600160a01b031633146103c75760405162461bcd60e51b81526004016102b890610747565b806002836040516103d8919061072b565b908152602001604051809103902090805190602001906103f9929190610474565b507fddc5a395ff29c22c0e109c1b1e032440d25c3f9452ffe7327b9dbb2f30fa632a828260405161042b929190610765565b60405180910390a15050565b508054610443906106f0565b6000825580601f10610453575050565b601f01602090049060005260206000209081019061047191906104f8565b50565b828054610480906106f0565b90600052602060002090601f0160209004810192826104a257600085556104e8565b82601f106104bb57805160ff19168380011785556104e8565b828001600101855582156104e8579182015b828111156104e85782518255916020019190600101906104cd565b506104f49291506104f8565b5090565b5b808211156104f457600081556001016104f9565b60005b83811015610528578181015183820152602001610510565b83811115610537576000848401525b50505050565b6000815180845261055581602086016020860161050d565b601f01601f19169290920160200192915050565b60208152600061057c602083018461053d565b9392505050565b634e487b7160e01b600052604160045260246000fd5b600082601f8301126105aa57600080fd5b813567ffffffffffffffff808211156105c5576105c5610583565b604051601f8301601f19908116603f011681019082821181831017156105ed576105ed610583565b8160405283815286602085880101111561060657600080fd5b836020870160208301376000602085830101528094505050505092915050565b60006020828403121561063857600080fd5b813567ffffffffffffffff81111561064f57600080fd5b61065b84828501610599565b949350505050565b60006020828403121561067557600080fd5b81356001600160a01b038116811461057c57600080fd5b6000806040838503121561069f57600080fd5b823567ffffffffffffffff808211156106b757600080fd5b6106c386838701610599565b935060208501359150808211156106d957600080fd5b506106e685828601610599565b9150509250929050565b600181811c9082168061070457607f821691505b6020821081141561072557634e487b7160e01b600052602260045260246000fd5b50919050565b6000825161073d81846020870161050d565b9190910192915050565b60208082526004908201526310b3b7bb60e11b604082015260600190565b604081526000610778604083018561053d565b828103602084015261078a818561053d565b9594505050505056fea26469706673582212201dce9e871bb9010286dfbbdb8a678ba2c02e634f10ca61c511eec4358b8e46be64736f6c63430008090033"

// DeployProof deploys a new Ethereum contract, binding an instance of Proof to it.
func DeployProof(auth *bind.TransactOpts, backend bind.ContractBackend, _proofName string, _govAddress common.Address) (common.Address, *types.Transaction, *Proof, error) {
	parsed, err := abi.JSON(strings.NewReader(ProofABI))
	if err != nil {
		return common.Address{}, nil, nil, err
	}

	address, tx, contract, err := bind.DeployContract(auth, parsed, common.FromHex(ProofBin), backend, _proofName, _govAddress)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &Proof{ProofCaller: ProofCaller{contract: contract}, ProofTransactor: ProofTransactor{contract: contract}, ProofFilterer: ProofFilterer{contract: contract}}, nil
}

// Proof is an auto generated Go binding around an Ethereum contract.
type Proof struct {
	ProofCaller     // Read-only binding to the contract
	ProofTransactor // Write-only binding to the contract
	ProofFilterer   // Log filterer for contract events
}

// ProofCaller is an auto generated read-only Go binding around an Ethereum contract.
type ProofCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ProofTransactor is an auto generated write-only Go binding around an Ethereum contract.
type ProofTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ProofFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type ProofFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ProofSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type ProofSession struct {
	Contract     *Proof            // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// ProofCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type ProofCallerSession struct {
	Contract *ProofCaller  // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts // Call options to use throughout this session
}

// ProofTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type ProofTransactorSession struct {
	Contract     *ProofTransactor  // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// ProofRaw is an auto generated low-level Go binding around an Ethereum contract.
type ProofRaw struct {
	Contract *Proof // Generic contract binding to access the raw methods on
}

// ProofCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type ProofCallerRaw struct {
	Contract *ProofCaller // Generic read-only contract binding to access the raw methods on
}

// ProofTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type ProofTransactorRaw struct {
	Contract *ProofTransactor // Generic write-only contract binding to access the raw methods on
}

// NewProof creates a new instance of Proof, bound to a specific deployed contract.
func NewProof(address common.Address, backend bind.ContractBackend) (*Proof, error) {
	contract, err := bindProof(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Proof{ProofCaller: ProofCaller{contract: contract}, ProofTransactor: ProofTransactor{contract: contract}, ProofFilterer: ProofFilterer{contract: contract}}, nil
}

// NewProofCaller creates a new read-only instance of Proof, bound to a specific deployed contract.
func NewProofCaller(address common.Address, caller bind.ContractCaller) (*ProofCaller, error) {
	contract, err := bindProof(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &ProofCaller{contract: contract}, nil
}

// NewProofTransactor creates a new write-only instance of Proof, bound to a specific deployed contract.
func NewProofTransactor(address common.Address, transactor bind.ContractTransactor) (*ProofTransactor, error) {
	contract, err := bindProof(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &ProofTransactor{contract: contract}, nil
}

// NewProofFilterer creates a new log filterer instance of Proof, bound to a specific deployed contract.
func NewProofFilterer(address common.Address, filterer bind.ContractFilterer) (*ProofFilterer, error) {
	contract, err := bindProof(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &ProofFilterer{contract: contract}, nil
}

// bindProof binds a generic wrapper to an already deployed contract.
func bindProof(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(ProofABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Proof *ProofRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Proof.Contract.ProofCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Proof *ProofRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Proof.Contract.ProofTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Proof *ProofRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Proof.Contract.ProofTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Proof *ProofCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Proof.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Proof *ProofTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Proof.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Proof *ProofTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Proof.Contract.contract.Transact(opts, method, params...)
}

// Get is a free data retrieval call binding the contract method 0x693ec85e.
//
// Solidity: function get(string _key) view returns(string)
func (_Proof *ProofCaller) Get(opts *bind.CallOpts, _key string) (string, error) {
	var out []interface{}
	err := _Proof.contract.Call(opts, &out, "get", _key)

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// Get is a free data retrieval call binding the contract method 0x693ec85e.
//
// Solidity: function get(string _key) view returns(string)
func (_Proof *ProofSession) Get(_key string) (string, error) {
	return _Proof.Contract.Get(&_Proof.CallOpts, _key)
}

// Get is a free data retrieval call binding the contract method 0x693ec85e.
//
// Solidity: function get(string _key) view returns(string)
func (_Proof *ProofCallerSession) Get(_key string) (string, error) {
	return _Proof.Contract.Get(&_Proof.CallOpts, _key)
}

// GovAddress is a free data retrieval call binding the contract method 0x46008a07.
//
// Solidity: function govAddress() view returns(address)
func (_Proof *ProofCaller) GovAddress(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Proof.contract.Call(opts, &out, "govAddress")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// GovAddress is a free data retrieval call binding the contract method 0x46008a07.
//
// Solidity: function govAddress() view returns(address)
func (_Proof *ProofSession) GovAddress() (common.Address, error) {
	return _Proof.Contract.GovAddress(&_Proof.CallOpts)
}

// GovAddress is a free data retrieval call binding the contract method 0x46008a07.
//
// Solidity: function govAddress() view returns(address)
func (_Proof *ProofCallerSession) GovAddress() (common.Address, error) {
	return _Proof.Contract.GovAddress(&_Proof.CallOpts)
}

// ProofName is a free data retrieval call binding the contract method 0x27d2fce5.
//
// Solidity: function proofName() view returns(string)
func (_Proof *ProofCaller) ProofName(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _Proof.contract.Call(opts, &out, "proofName")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// ProofName is a free data retrieval call binding the contract method 0x27d2fce5.
//
// Solidity: function proofName() view returns(string)
func (_Proof *ProofSession) ProofName() (string, error) {
	return _Proof.Contract.ProofName(&_Proof.CallOpts)
}

// ProofName is a free data retrieval call binding the contract method 0x27d2fce5.
//
// Solidity: function proofName() view returns(string)
func (_Proof *ProofCallerSession) ProofName() (string, error) {
	return _Proof.Contract.ProofName(&_Proof.CallOpts)
}

// Proofs is a free data retrieval call binding the contract method 0x3b18f691.
//
// Solidity: function proofs(string ) view returns(string)
func (_Proof *ProofCaller) Proofs(opts *bind.CallOpts, arg0 string) (string, error) {
	var out []interface{}
	err := _Proof.contract.Call(opts, &out, "proofs", arg0)

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// Proofs is a free data retrieval call binding the contract method 0x3b18f691.
//
// Solidity: function proofs(string ) view returns(string)
func (_Proof *ProofSession) Proofs(arg0 string) (string, error) {
	return _Proof.Contract.Proofs(&_Proof.CallOpts, arg0)
}

// Proofs is a free data retrieval call binding the contract method 0x3b18f691.
//
// Solidity: function proofs(string ) view returns(string)
func (_Proof *ProofCallerSession) Proofs(arg0 string) (string, error) {
	return _Proof.Contract.Proofs(&_Proof.CallOpts, arg0)
}

// Remove is a paid mutator transaction binding the contract method 0x80599e4b.
//
// Solidity: function remove(string _key) returns()
func (_Proof *ProofTransactor) Remove(opts *bind.TransactOpts, _key string) (*types.Transaction, error) {
	return _Proof.contract.Transact(opts, "remove", _key)
}

// Remove is a paid mutator transaction binding the contract method 0x80599e4b.
//
// Solidity: function remove(string _key) returns()
func (_Proof *ProofSession) Remove(_key string) (*types.Transaction, error) {
	return _Proof.Contract.Remove(&_Proof.TransactOpts, _key)
}

// Remove is a paid mutator transaction binding the contract method 0x80599e4b.
//
// Solidity: function remove(string _key) returns()
func (_Proof *ProofTransactorSession) Remove(_key string) (*types.Transaction, error) {
	return _Proof.Contract.Remove(&_Proof.TransactOpts, _key)
}

// Set is a paid mutator transaction binding the contract method 0xe942b516.
//
// Solidity: function set(string _key, string _value) returns()
func (_Proof *ProofTransactor) Set(opts *bind.TransactOpts, _key string, _value string) (*types.Transaction, error) {
	return _Proof.contract.Transact(opts, "set", _key, _value)
}

// Set is a paid mutator transaction binding the contract method 0xe942b516.
//
// Solidity: function set(string _key, string _value) returns()
func (_Proof *ProofSession) Set(_key string, _value string) (*types.Transaction, error) {
	return _Proof.Contract.Set(&_Proof.TransactOpts, _key, _value)
}

// Set is a paid mutator transaction binding the contract method 0xe942b516.
//
// Solidity: function set(string _key, string _value) returns()
func (_Proof *ProofTransactorSession) Set(_key string, _value string) (*types.Transaction, error) {
	return _Proof.Contract.Set(&_Proof.TransactOpts, _key, _value)
}

// SetGov is a paid mutator transaction binding the contract method 0xcfad57a2.
//
// Solidity: function setGov(address _govAddress) returns()
func (_Proof *ProofTransactor) SetGov(opts *bind.TransactOpts, _govAddress common.Address) (*types.Transaction, error) {
	return _Proof.contract.Transact(opts, "setGov", _govAddress)
}

// SetGov is a paid mutator transaction binding the contract method 0xcfad57a2.
//
// Solidity: function setGov(address _govAddress) returns()
func (_Proof *ProofSession) SetGov(_govAddress common.Address) (*types.Transaction, error) {
	return _Proof.Contract.SetGov(&_Proof.TransactOpts, _govAddress)
}

// SetGov is a paid mutator transaction binding the contract method 0xcfad57a2.
//
// Solidity: function setGov(address _govAddress) returns()
func (_Proof *ProofTransactorSession) SetGov(_govAddress common.Address) (*types.Transaction, error) {
	return _Proof.Contract.SetGov(&_Proof.TransactOpts, _govAddress)
}

// ProofRemoveIterator is returned from FilterRemove and is used to iterate over the raw logs and unpacked data for Remove events raised by the Proof contract.
type ProofRemoveIterator struct {
	Event *ProofRemove // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *ProofRemoveIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ProofRemove)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(ProofRemove)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *ProofRemoveIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ProofRemoveIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ProofRemove represents a Remove event raised by the Proof contract.
type ProofRemove struct {
	Key string
	Raw types.Log // Blockchain specific contextual infos
}

// FilterRemove is a free log retrieval operation binding the contract event 0x834a2d47e948021d7136fb7275b3f1e1feae6333c0d683e8c13f901667defd8c.
//
// Solidity: event Remove(string _key)
func (_Proof *ProofFilterer) FilterRemove(opts *bind.FilterOpts) (*ProofRemoveIterator, error) {

	logs, sub, err := _Proof.contract.FilterLogs(opts, "Remove")
	if err != nil {
		return nil, err
	}
	return &ProofRemoveIterator{contract: _Proof.contract, event: "Remove", logs: logs, sub: sub}, nil
}

// WatchRemove is a free log subscription operation binding the contract event 0x834a2d47e948021d7136fb7275b3f1e1feae6333c0d683e8c13f901667defd8c.
//
// Solidity: event Remove(string _key)
func (_Proof *ProofFilterer) WatchRemove(opts *bind.WatchOpts, sink chan<- *ProofRemove) (event.Subscription, error) {

	logs, sub, err := _Proof.contract.WatchLogs(opts, "Remove")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ProofRemove)
				if err := _Proof.contract.UnpackLog(event, "Remove", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseRemove is a log parse operation binding the contract event 0x834a2d47e948021d7136fb7275b3f1e1feae6333c0d683e8c13f901667defd8c.
//
// Solidity: event Remove(string _key)
func (_Proof *ProofFilterer) ParseRemove(log types.Log) (*ProofRemove, error) {
	event := new(ProofRemove)
	if err := _Proof.contract.UnpackLog(event, "Remove", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ProofSetIterator is returned from FilterSet and is used to iterate over the raw logs and unpacked data for Set events raised by the Proof contract.
type ProofSetIterator struct {
	Event *ProofSet // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *ProofSetIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ProofSet)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(ProofSet)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *ProofSetIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ProofSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ProofSet represents a Set event raised by the Proof contract.
type ProofSet struct {
	Key   string
	Value string
	Raw   types.Log // Blockchain specific contextual infos
}

// FilterSet is a free log retrieval operation binding the contract event 0xddc5a395ff29c22c0e109c1b1e032440d25c3f9452ffe7327b9dbb2f30fa632a.
//
// Solidity: event Set(string _key, string _value)
func (_Proof *ProofFilterer) FilterSet(opts *bind.FilterOpts) (*ProofSetIterator, error) {

	logs, sub, err := _Proof.contract.FilterLogs(opts, "Set")
	if err != nil {
		return nil, err
	}
	return &ProofSetIterator{contract: _Proof.contract, event: "Set", logs: logs, sub: sub}, nil
}

// WatchSet is a free log subscription operation binding the contract event 0xddc5a395ff29c22c0e109c1b1e032440d25c3f9452ffe7327b9dbb2f30fa632a.
//
// Solidity: event Set(string _key, string _value)
func (_Proof *ProofFilterer) WatchSet(opts *bind.WatchOpts, sink chan<- *ProofSet) (event.Subscription, error) {

	logs, sub, err := _Proof.contract.WatchLogs(opts, "Set")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ProofSet)
				if err := _Proof.contract.UnpackLog(event, "Set", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseSet is a log parse operation binding the contract event 0xddc5a395ff29c22c0e109c1b1e032440d25c3f9452ffe7327b9dbb2f30fa632a.
//
// Solidity: event Set(string _key, string _value)
func (_Proof *ProofFilterer) ParseSet(log types.Log) (*ProofSet, error) {
	event := new(ProofSet)
	if err := _Proof.contract.UnpackLog(event, "Set", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ProofSetGovIterator is returned from FilterSetGov and is used to iterate over the raw logs and unpacked data for SetGov events raised by the Proof contract.
type ProofSetGovIterator struct {
	Event *ProofSetGov // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *ProofSetGovIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ProofSetGov)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(ProofSetGov)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *ProofSetGovIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ProofSetGovIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ProofSetGov represents a SetGov event raised by the Proof contract.
type ProofSetGov struct {
	GovAddress common.Address
	Raw        types.Log // Blockchain specific contextual infos
}

// FilterSetGov is a free log retrieval operation binding the contract event 0x91a8c1cc2d4a3bb60738481947a00cbb9899c822916694cf8bb1d68172fdcd54.
//
// Solidity: event SetGov(address _govAddress)
func (_Proof *ProofFilterer) FilterSetGov(opts *bind.FilterOpts) (*ProofSetGovIterator, error) {

	logs, sub, err := _Proof.contract.FilterLogs(opts, "SetGov")
	if err != nil {
		return nil, err
	}
	return &ProofSetGovIterator{contract: _Proof.contract, event: "SetGov", logs: logs, sub: sub}, nil
}

// WatchSetGov is a free log subscription operation binding the contract event 0x91a8c1cc2d4a3bb60738481947a00cbb9899c822916694cf8bb1d68172fdcd54.
//
// Solidity: event SetGov(address _govAddress)
func (_Proof *ProofFilterer) WatchSetGov(opts *bind.WatchOpts, sink chan<- *ProofSetGov) (event.Subscription, error) {

	logs, sub, err := _Proof.contract.WatchLogs(opts, "SetGov")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ProofSetGov)
				if err := _Proof.contract.UnpackLog(event, "SetGov", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseSetGov is a log parse operation binding the contract event 0x91a8c1cc2d4a3bb60738481947a00cbb9899c822916694cf8bb1d68172fdcd54.
//
// Solidity: event SetGov(address _govAddress)
func (_Proof *ProofFilterer) ParseSetGov(log types.Log) (*ProofSetGov, error) {
	event := new(ProofSetGov)
	if err := _Proof.contract.UnpackLog(event, "SetGov", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
