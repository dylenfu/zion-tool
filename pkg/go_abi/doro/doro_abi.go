// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package doro

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
	MethodSetDoro = "setDoro"

	MethodData = "data"
)

// DoroABI is the input ABI used to generate the binding from.
const DoroABI = "[{\"constant\":true,\"inputs\":[],\"name\":\"data\",\"outputs\":[{\"internalType\":\"uint64\",\"name\":\"\",\"type\":\"uint64\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"num\",\"type\":\"uint64\"}],\"name\":\"setDoro\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]"

// DoroFuncSigs maps the 4-byte function signature to its string representation.
var DoroFuncSigs = map[string]string{
	"73d4a13a": "data()",
	"517b2ea7": "setDoro(uint64)",
}

// DoroBin is the compiled bytecode used for deploying new contracts.
var DoroBin = "0x608060405234801561001057600080fd5b5060fd8061001f6000396000f3fe6080604052348015600f57600080fd5b506004361060325760003560e01c8063517b2ea714603757806373d4a13a14606f575b600080fd5b605b60048036036020811015604b57600080fd5b503567ffffffffffffffff166092565b604080519115158252519081900360200190f35b607560b8565b6040805167ffffffffffffffff9092168252519081900360200190f35b6000805467ffffffffffffffff831667ffffffffffffffff199091161790556001919050565b60005467ffffffffffffffff168156fea265627a7a723158202f5d331dfdec468c3a849101fc87bdbb580f207958159fc9dc8798be556cf55464736f6c63430005110032"

// DeployDoro deploys a new Ethereum contract, binding an instance of Doro to it.
func DeployDoro(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *Doro, error) {
	parsed, err := abi.JSON(strings.NewReader(DoroABI))
	if err != nil {
		return common.Address{}, nil, nil, err
	}

	address, tx, contract, err := bind.DeployContract(auth, parsed, common.FromHex(DoroBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &Doro{DoroCaller: DoroCaller{contract: contract}, DoroTransactor: DoroTransactor{contract: contract}, DoroFilterer: DoroFilterer{contract: contract}}, nil
}

// Doro is an auto generated Go binding around an Ethereum contract.
type Doro struct {
	DoroCaller     // Read-only binding to the contract
	DoroTransactor // Write-only binding to the contract
	DoroFilterer   // Log filterer for contract events
}

// DoroCaller is an auto generated read-only Go binding around an Ethereum contract.
type DoroCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// DoroTransactor is an auto generated write-only Go binding around an Ethereum contract.
type DoroTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// DoroFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type DoroFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// DoroSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type DoroSession struct {
	Contract     *Doro             // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// DoroCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type DoroCallerSession struct {
	Contract *DoroCaller   // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts // Call options to use throughout this session
}

// DoroTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type DoroTransactorSession struct {
	Contract     *DoroTransactor   // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// DoroRaw is an auto generated low-level Go binding around an Ethereum contract.
type DoroRaw struct {
	Contract *Doro // Generic contract binding to access the raw methods on
}

// DoroCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type DoroCallerRaw struct {
	Contract *DoroCaller // Generic read-only contract binding to access the raw methods on
}

// DoroTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type DoroTransactorRaw struct {
	Contract *DoroTransactor // Generic write-only contract binding to access the raw methods on
}

// NewDoro creates a new instance of Doro, bound to a specific deployed contract.
func NewDoro(address common.Address, backend bind.ContractBackend) (*Doro, error) {
	contract, err := bindDoro(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Doro{DoroCaller: DoroCaller{contract: contract}, DoroTransactor: DoroTransactor{contract: contract}, DoroFilterer: DoroFilterer{contract: contract}}, nil
}

// NewDoroCaller creates a new read-only instance of Doro, bound to a specific deployed contract.
func NewDoroCaller(address common.Address, caller bind.ContractCaller) (*DoroCaller, error) {
	contract, err := bindDoro(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &DoroCaller{contract: contract}, nil
}

// NewDoroTransactor creates a new write-only instance of Doro, bound to a specific deployed contract.
func NewDoroTransactor(address common.Address, transactor bind.ContractTransactor) (*DoroTransactor, error) {
	contract, err := bindDoro(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &DoroTransactor{contract: contract}, nil
}

// NewDoroFilterer creates a new log filterer instance of Doro, bound to a specific deployed contract.
func NewDoroFilterer(address common.Address, filterer bind.ContractFilterer) (*DoroFilterer, error) {
	contract, err := bindDoro(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &DoroFilterer{contract: contract}, nil
}

// bindDoro binds a generic wrapper to an already deployed contract.
func bindDoro(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(DoroABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Doro *DoroRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Doro.Contract.DoroCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Doro *DoroRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Doro.Contract.DoroTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Doro *DoroRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Doro.Contract.DoroTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Doro *DoroCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Doro.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Doro *DoroTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Doro.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Doro *DoroTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Doro.Contract.contract.Transact(opts, method, params...)
}

// Data is a free data retrieval call binding the contract method 0x73d4a13a.
//
// Solidity: function data() view returns(uint64)
func (_Doro *DoroCaller) Data(opts *bind.CallOpts) (uint64, error) {
	var out []interface{}
	err := _Doro.contract.Call(opts, &out, "data")

	if err != nil {
		return *new(uint64), err
	}

	out0 := *abi.ConvertType(out[0], new(uint64)).(*uint64)

	return out0, err

}

// Data is a free data retrieval call binding the contract method 0x73d4a13a.
//
// Solidity: function data() view returns(uint64)
func (_Doro *DoroSession) Data() (uint64, error) {
	return _Doro.Contract.Data(&_Doro.CallOpts)
}

// Data is a free data retrieval call binding the contract method 0x73d4a13a.
//
// Solidity: function data() view returns(uint64)
func (_Doro *DoroCallerSession) Data() (uint64, error) {
	return _Doro.Contract.Data(&_Doro.CallOpts)
}

// SetDoro is a paid mutator transaction binding the contract method 0x517b2ea7.
//
// Solidity: function setDoro(uint64 num) returns(bool)
func (_Doro *DoroTransactor) SetDoro(opts *bind.TransactOpts, num uint64) (*types.Transaction, error) {
	return _Doro.contract.Transact(opts, "setDoro", num)
}

// SetDoro is a paid mutator transaction binding the contract method 0x517b2ea7.
//
// Solidity: function setDoro(uint64 num) returns(bool)
func (_Doro *DoroSession) SetDoro(num uint64) (*types.Transaction, error) {
	return _Doro.Contract.SetDoro(&_Doro.TransactOpts, num)
}

// SetDoro is a paid mutator transaction binding the contract method 0x517b2ea7.
//
// Solidity: function setDoro(uint64 num) returns(bool)
func (_Doro *DoroTransactorSession) SetDoro(num uint64) (*types.Transaction, error) {
	return _Doro.Contract.SetDoro(&_Doro.TransactOpts, num)
}

