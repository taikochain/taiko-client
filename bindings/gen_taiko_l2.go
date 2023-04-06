// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package bindings

import (
	"errors"
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
	_ = errors.New
	_ = big.NewInt
	_ = strings.NewReader
	_ = ethereum.NotFound
	_ = bind.Bind
	_ = common.Big1
	_ = types.BloomLookup
	_ = event.NewSubscription
	_ = abi.ConvertType
)

// TaikoL2EIP1559Params is an auto generated low-level Go binding around an user-defined struct.
type TaikoL2EIP1559Params struct {
	Basefee            uint64
	GasIssuedPerSecond uint64
	GasExcessMax       uint64
	GasTarget          uint64
	Ratio2x1x          uint64
}

// TaikoL2ClientMetaData contains all meta data concerning the TaikoL2Client contract.
var TaikoL2ClientMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"expected\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"actual\",\"type\":\"uint64\"}],\"name\":\"L2_BASEFEE_MISMATCH\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"L2_INVALID_1559_PARAMS\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"L2_INVALID_CHAIN_ID\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"L2_INVALID_GOLDEN_TOUCH_K\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"L2_INVALID_SENDER\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"expected\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"actual\",\"type\":\"bytes32\"}],\"name\":\"L2_PUBLIC_INPUT_HASH_MISMATCH\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"L2_TOO_LATE\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"M1559_OUT_OF_STOCK\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"M1559_OUT_OF_STOCK\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"expected\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"actual\",\"type\":\"uint64\"}],\"name\":\"M1559_UNEXPECTED_CHANGE\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"expected\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"actual\",\"type\":\"uint64\"}],\"name\":\"M1559_UNEXPECTED_CHANGE\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"Overflow\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"RESOLVER_DENIED\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"RESOLVER_INVALID_ADDR\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"number\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"basefee\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"gaslimit\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"timestamp\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"parentHash\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"prevrandao\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"coinbase\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"chainid\",\"type\":\"uint32\"}],\"name\":\"BlockVars\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint8\",\"name\":\"version\",\"type\":\"uint8\"}],\"name\":\"Initialized\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"previousOwner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"srcHeight\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"blockHash\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"signalRoot\",\"type\":\"bytes32\"}],\"name\":\"XchainSynced\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"ANCHOR_GAS_COST\",\"outputs\":[{\"internalType\":\"uint64\",\"name\":\"\",\"type\":\"uint64\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"GOLDEN_TOUCH_ADDRESS\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"GOLDEN_TOUCH_PRIVATEKEY\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"__reserved1\",\"outputs\":[{\"internalType\":\"uint64\",\"name\":\"\",\"type\":\"uint64\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"addressManager\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"l1Hash\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"l1SignalRoot\",\"type\":\"bytes32\"},{\"internalType\":\"uint64\",\"name\":\"l1Height\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"parentGasUsed\",\"type\":\"uint64\"}],\"name\":\"anchor\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"gasExcess\",\"outputs\":[{\"internalType\":\"uint64\",\"name\":\"\",\"type\":\"uint64\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"gasIssuedPerSecond\",\"outputs\":[{\"internalType\":\"uint64\",\"name\":\"\",\"type\":\"uint64\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint32\",\"name\":\"timeSinceNow\",\"type\":\"uint32\"},{\"internalType\":\"uint64\",\"name\":\"gasLimit\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"parentGasUsed\",\"type\":\"uint64\"}],\"name\":\"getBasefee\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"_basefee\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"number\",\"type\":\"uint256\"}],\"name\":\"getBlockHash\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"number\",\"type\":\"uint256\"}],\"name\":\"getXchainBlockHash\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"number\",\"type\":\"uint256\"}],\"name\":\"getXchainSignalRoot\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_addressManager\",\"type\":\"address\"},{\"components\":[{\"internalType\":\"uint64\",\"name\":\"basefee\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"gasIssuedPerSecond\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"gasExcessMax\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"gasTarget\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"ratio2x1x\",\"type\":\"uint64\"}],\"internalType\":\"structTaikoL2.EIP1559Params\",\"name\":\"_param1559\",\"type\":\"tuple\"}],\"name\":\"init\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"chainId\",\"type\":\"uint256\"},{\"internalType\":\"string\",\"name\":\"name\",\"type\":\"string\"}],\"name\":\"keyForName\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"latestSyncedL1Height\",\"outputs\":[{\"internalType\":\"uint64\",\"name\":\"\",\"type\":\"uint64\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"parentTimestamp\",\"outputs\":[{\"internalType\":\"uint64\",\"name\":\"\",\"type\":\"uint64\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"publicInputHash\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"renounceOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"name\",\"type\":\"string\"},{\"internalType\":\"bool\",\"name\":\"allowZeroAddress\",\"type\":\"bool\"}],\"name\":\"resolve\",\"outputs\":[{\"internalType\":\"addresspayable\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"chainId\",\"type\":\"uint256\"},{\"internalType\":\"string\",\"name\":\"name\",\"type\":\"string\"},{\"internalType\":\"bool\",\"name\":\"allowZeroAddress\",\"type\":\"bool\"}],\"name\":\"resolve\",\"outputs\":[{\"internalType\":\"addresspayable\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"digest\",\"type\":\"bytes32\"},{\"internalType\":\"uint8\",\"name\":\"k\",\"type\":\"uint8\"}],\"name\":\"signAnchor\",\"outputs\":[{\"internalType\":\"uint8\",\"name\":\"v\",\"type\":\"uint8\"},{\"internalType\":\"uint256\",\"name\":\"r\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"s\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"xscale\",\"outputs\":[{\"internalType\":\"uint64\",\"name\":\"\",\"type\":\"uint64\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"yscale\",\"outputs\":[{\"internalType\":\"uint128\",\"name\":\"\",\"type\":\"uint128\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
}

// TaikoL2ClientABI is the input ABI used to generate the binding from.
// Deprecated: Use TaikoL2ClientMetaData.ABI instead.
var TaikoL2ClientABI = TaikoL2ClientMetaData.ABI

// TaikoL2Client is an auto generated Go binding around an Ethereum contract.
type TaikoL2Client struct {
	TaikoL2ClientCaller     // Read-only binding to the contract
	TaikoL2ClientTransactor // Write-only binding to the contract
	TaikoL2ClientFilterer   // Log filterer for contract events
}

// TaikoL2ClientCaller is an auto generated read-only Go binding around an Ethereum contract.
type TaikoL2ClientCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// TaikoL2ClientTransactor is an auto generated write-only Go binding around an Ethereum contract.
type TaikoL2ClientTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// TaikoL2ClientFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type TaikoL2ClientFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// TaikoL2ClientSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type TaikoL2ClientSession struct {
	Contract     *TaikoL2Client    // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// TaikoL2ClientCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type TaikoL2ClientCallerSession struct {
	Contract *TaikoL2ClientCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts        // Call options to use throughout this session
}

// TaikoL2ClientTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type TaikoL2ClientTransactorSession struct {
	Contract     *TaikoL2ClientTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts        // Transaction auth options to use throughout this session
}

// TaikoL2ClientRaw is an auto generated low-level Go binding around an Ethereum contract.
type TaikoL2ClientRaw struct {
	Contract *TaikoL2Client // Generic contract binding to access the raw methods on
}

// TaikoL2ClientCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type TaikoL2ClientCallerRaw struct {
	Contract *TaikoL2ClientCaller // Generic read-only contract binding to access the raw methods on
}

// TaikoL2ClientTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type TaikoL2ClientTransactorRaw struct {
	Contract *TaikoL2ClientTransactor // Generic write-only contract binding to access the raw methods on
}

// NewTaikoL2Client creates a new instance of TaikoL2Client, bound to a specific deployed contract.
func NewTaikoL2Client(address common.Address, backend bind.ContractBackend) (*TaikoL2Client, error) {
	contract, err := bindTaikoL2Client(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &TaikoL2Client{TaikoL2ClientCaller: TaikoL2ClientCaller{contract: contract}, TaikoL2ClientTransactor: TaikoL2ClientTransactor{contract: contract}, TaikoL2ClientFilterer: TaikoL2ClientFilterer{contract: contract}}, nil
}

// NewTaikoL2ClientCaller creates a new read-only instance of TaikoL2Client, bound to a specific deployed contract.
func NewTaikoL2ClientCaller(address common.Address, caller bind.ContractCaller) (*TaikoL2ClientCaller, error) {
	contract, err := bindTaikoL2Client(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &TaikoL2ClientCaller{contract: contract}, nil
}

// NewTaikoL2ClientTransactor creates a new write-only instance of TaikoL2Client, bound to a specific deployed contract.
func NewTaikoL2ClientTransactor(address common.Address, transactor bind.ContractTransactor) (*TaikoL2ClientTransactor, error) {
	contract, err := bindTaikoL2Client(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &TaikoL2ClientTransactor{contract: contract}, nil
}

// NewTaikoL2ClientFilterer creates a new log filterer instance of TaikoL2Client, bound to a specific deployed contract.
func NewTaikoL2ClientFilterer(address common.Address, filterer bind.ContractFilterer) (*TaikoL2ClientFilterer, error) {
	contract, err := bindTaikoL2Client(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &TaikoL2ClientFilterer{contract: contract}, nil
}

// bindTaikoL2Client binds a generic wrapper to an already deployed contract.
func bindTaikoL2Client(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := TaikoL2ClientMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_TaikoL2Client *TaikoL2ClientRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _TaikoL2Client.Contract.TaikoL2ClientCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_TaikoL2Client *TaikoL2ClientRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _TaikoL2Client.Contract.TaikoL2ClientTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_TaikoL2Client *TaikoL2ClientRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _TaikoL2Client.Contract.TaikoL2ClientTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_TaikoL2Client *TaikoL2ClientCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _TaikoL2Client.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_TaikoL2Client *TaikoL2ClientTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _TaikoL2Client.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_TaikoL2Client *TaikoL2ClientTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _TaikoL2Client.Contract.contract.Transact(opts, method, params...)
}

// ANCHORGASCOST is a free data retrieval call binding the contract method 0x53ecb5b6.
//
// Solidity: function ANCHOR_GAS_COST() view returns(uint64)
func (_TaikoL2Client *TaikoL2ClientCaller) ANCHORGASCOST(opts *bind.CallOpts) (uint64, error) {
	var out []interface{}
	err := _TaikoL2Client.contract.Call(opts, &out, "ANCHOR_GAS_COST")

	if err != nil {
		return *new(uint64), err
	}

	out0 := *abi.ConvertType(out[0], new(uint64)).(*uint64)

	return out0, err

}

// ANCHORGASCOST is a free data retrieval call binding the contract method 0x53ecb5b6.
//
// Solidity: function ANCHOR_GAS_COST() view returns(uint64)
func (_TaikoL2Client *TaikoL2ClientSession) ANCHORGASCOST() (uint64, error) {
	return _TaikoL2Client.Contract.ANCHORGASCOST(&_TaikoL2Client.CallOpts)
}

// ANCHORGASCOST is a free data retrieval call binding the contract method 0x53ecb5b6.
//
// Solidity: function ANCHOR_GAS_COST() view returns(uint64)
func (_TaikoL2Client *TaikoL2ClientCallerSession) ANCHORGASCOST() (uint64, error) {
	return _TaikoL2Client.Contract.ANCHORGASCOST(&_TaikoL2Client.CallOpts)
}

// GOLDENTOUCHADDRESS is a free data retrieval call binding the contract method 0x9ee512f2.
//
// Solidity: function GOLDEN_TOUCH_ADDRESS() view returns(address)
func (_TaikoL2Client *TaikoL2ClientCaller) GOLDENTOUCHADDRESS(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _TaikoL2Client.contract.Call(opts, &out, "GOLDEN_TOUCH_ADDRESS")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// GOLDENTOUCHADDRESS is a free data retrieval call binding the contract method 0x9ee512f2.
//
// Solidity: function GOLDEN_TOUCH_ADDRESS() view returns(address)
func (_TaikoL2Client *TaikoL2ClientSession) GOLDENTOUCHADDRESS() (common.Address, error) {
	return _TaikoL2Client.Contract.GOLDENTOUCHADDRESS(&_TaikoL2Client.CallOpts)
}

// GOLDENTOUCHADDRESS is a free data retrieval call binding the contract method 0x9ee512f2.
//
// Solidity: function GOLDEN_TOUCH_ADDRESS() view returns(address)
func (_TaikoL2Client *TaikoL2ClientCallerSession) GOLDENTOUCHADDRESS() (common.Address, error) {
	return _TaikoL2Client.Contract.GOLDENTOUCHADDRESS(&_TaikoL2Client.CallOpts)
}

// GOLDENTOUCHPRIVATEKEY is a free data retrieval call binding the contract method 0x10da3738.
//
// Solidity: function GOLDEN_TOUCH_PRIVATEKEY() view returns(uint256)
func (_TaikoL2Client *TaikoL2ClientCaller) GOLDENTOUCHPRIVATEKEY(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _TaikoL2Client.contract.Call(opts, &out, "GOLDEN_TOUCH_PRIVATEKEY")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GOLDENTOUCHPRIVATEKEY is a free data retrieval call binding the contract method 0x10da3738.
//
// Solidity: function GOLDEN_TOUCH_PRIVATEKEY() view returns(uint256)
func (_TaikoL2Client *TaikoL2ClientSession) GOLDENTOUCHPRIVATEKEY() (*big.Int, error) {
	return _TaikoL2Client.Contract.GOLDENTOUCHPRIVATEKEY(&_TaikoL2Client.CallOpts)
}

// GOLDENTOUCHPRIVATEKEY is a free data retrieval call binding the contract method 0x10da3738.
//
// Solidity: function GOLDEN_TOUCH_PRIVATEKEY() view returns(uint256)
func (_TaikoL2Client *TaikoL2ClientCallerSession) GOLDENTOUCHPRIVATEKEY() (*big.Int, error) {
	return _TaikoL2Client.Contract.GOLDENTOUCHPRIVATEKEY(&_TaikoL2Client.CallOpts)
}

// Reserved1 is a free data retrieval call binding the contract method 0xf5d11edc.
//
// Solidity: function __reserved1() view returns(uint64)
func (_TaikoL2Client *TaikoL2ClientCaller) Reserved1(opts *bind.CallOpts) (uint64, error) {
	var out []interface{}
	err := _TaikoL2Client.contract.Call(opts, &out, "__reserved1")

	if err != nil {
		return *new(uint64), err
	}

	out0 := *abi.ConvertType(out[0], new(uint64)).(*uint64)

	return out0, err

}

// Reserved1 is a free data retrieval call binding the contract method 0xf5d11edc.
//
// Solidity: function __reserved1() view returns(uint64)
func (_TaikoL2Client *TaikoL2ClientSession) Reserved1() (uint64, error) {
	return _TaikoL2Client.Contract.Reserved1(&_TaikoL2Client.CallOpts)
}

// Reserved1 is a free data retrieval call binding the contract method 0xf5d11edc.
//
// Solidity: function __reserved1() view returns(uint64)
func (_TaikoL2Client *TaikoL2ClientCallerSession) Reserved1() (uint64, error) {
	return _TaikoL2Client.Contract.Reserved1(&_TaikoL2Client.CallOpts)
}

// AddressManager is a free data retrieval call binding the contract method 0x3ab76e9f.
//
// Solidity: function addressManager() view returns(address)
func (_TaikoL2Client *TaikoL2ClientCaller) AddressManager(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _TaikoL2Client.contract.Call(opts, &out, "addressManager")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// AddressManager is a free data retrieval call binding the contract method 0x3ab76e9f.
//
// Solidity: function addressManager() view returns(address)
func (_TaikoL2Client *TaikoL2ClientSession) AddressManager() (common.Address, error) {
	return _TaikoL2Client.Contract.AddressManager(&_TaikoL2Client.CallOpts)
}

// AddressManager is a free data retrieval call binding the contract method 0x3ab76e9f.
//
// Solidity: function addressManager() view returns(address)
func (_TaikoL2Client *TaikoL2ClientCallerSession) AddressManager() (common.Address, error) {
	return _TaikoL2Client.Contract.AddressManager(&_TaikoL2Client.CallOpts)
}

// GasExcess is a free data retrieval call binding the contract method 0xf535bd56.
//
// Solidity: function gasExcess() view returns(uint64)
func (_TaikoL2Client *TaikoL2ClientCaller) GasExcess(opts *bind.CallOpts) (uint64, error) {
	var out []interface{}
	err := _TaikoL2Client.contract.Call(opts, &out, "gasExcess")

	if err != nil {
		return *new(uint64), err
	}

	out0 := *abi.ConvertType(out[0], new(uint64)).(*uint64)

	return out0, err

}

// GasExcess is a free data retrieval call binding the contract method 0xf535bd56.
//
// Solidity: function gasExcess() view returns(uint64)
func (_TaikoL2Client *TaikoL2ClientSession) GasExcess() (uint64, error) {
	return _TaikoL2Client.Contract.GasExcess(&_TaikoL2Client.CallOpts)
}

// GasExcess is a free data retrieval call binding the contract method 0xf535bd56.
//
// Solidity: function gasExcess() view returns(uint64)
func (_TaikoL2Client *TaikoL2ClientCallerSession) GasExcess() (uint64, error) {
	return _TaikoL2Client.Contract.GasExcess(&_TaikoL2Client.CallOpts)
}

// GasIssuedPerSecond is a free data retrieval call binding the contract method 0x210c9fe8.
//
// Solidity: function gasIssuedPerSecond() view returns(uint64)
func (_TaikoL2Client *TaikoL2ClientCaller) GasIssuedPerSecond(opts *bind.CallOpts) (uint64, error) {
	var out []interface{}
	err := _TaikoL2Client.contract.Call(opts, &out, "gasIssuedPerSecond")

	if err != nil {
		return *new(uint64), err
	}

	out0 := *abi.ConvertType(out[0], new(uint64)).(*uint64)

	return out0, err

}

// GasIssuedPerSecond is a free data retrieval call binding the contract method 0x210c9fe8.
//
// Solidity: function gasIssuedPerSecond() view returns(uint64)
func (_TaikoL2Client *TaikoL2ClientSession) GasIssuedPerSecond() (uint64, error) {
	return _TaikoL2Client.Contract.GasIssuedPerSecond(&_TaikoL2Client.CallOpts)
}

// GasIssuedPerSecond is a free data retrieval call binding the contract method 0x210c9fe8.
//
// Solidity: function gasIssuedPerSecond() view returns(uint64)
func (_TaikoL2Client *TaikoL2ClientCallerSession) GasIssuedPerSecond() (uint64, error) {
	return _TaikoL2Client.Contract.GasIssuedPerSecond(&_TaikoL2Client.CallOpts)
}

// GetBasefee is a free data retrieval call binding the contract method 0xe1848cb0.
//
// Solidity: function getBasefee(uint32 timeSinceNow, uint64 gasLimit, uint64 parentGasUsed) view returns(uint256 _basefee)
func (_TaikoL2Client *TaikoL2ClientCaller) GetBasefee(opts *bind.CallOpts, timeSinceNow uint32, gasLimit uint64, parentGasUsed uint64) (*big.Int, error) {
	var out []interface{}
	err := _TaikoL2Client.contract.Call(opts, &out, "getBasefee", timeSinceNow, gasLimit, parentGasUsed)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetBasefee is a free data retrieval call binding the contract method 0xe1848cb0.
//
// Solidity: function getBasefee(uint32 timeSinceNow, uint64 gasLimit, uint64 parentGasUsed) view returns(uint256 _basefee)
func (_TaikoL2Client *TaikoL2ClientSession) GetBasefee(timeSinceNow uint32, gasLimit uint64, parentGasUsed uint64) (*big.Int, error) {
	return _TaikoL2Client.Contract.GetBasefee(&_TaikoL2Client.CallOpts, timeSinceNow, gasLimit, parentGasUsed)
}

// GetBasefee is a free data retrieval call binding the contract method 0xe1848cb0.
//
// Solidity: function getBasefee(uint32 timeSinceNow, uint64 gasLimit, uint64 parentGasUsed) view returns(uint256 _basefee)
func (_TaikoL2Client *TaikoL2ClientCallerSession) GetBasefee(timeSinceNow uint32, gasLimit uint64, parentGasUsed uint64) (*big.Int, error) {
	return _TaikoL2Client.Contract.GetBasefee(&_TaikoL2Client.CallOpts, timeSinceNow, gasLimit, parentGasUsed)
}

// GetBlockHash is a free data retrieval call binding the contract method 0xee82ac5e.
//
// Solidity: function getBlockHash(uint256 number) view returns(bytes32)
func (_TaikoL2Client *TaikoL2ClientCaller) GetBlockHash(opts *bind.CallOpts, number *big.Int) ([32]byte, error) {
	var out []interface{}
	err := _TaikoL2Client.contract.Call(opts, &out, "getBlockHash", number)

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// GetBlockHash is a free data retrieval call binding the contract method 0xee82ac5e.
//
// Solidity: function getBlockHash(uint256 number) view returns(bytes32)
func (_TaikoL2Client *TaikoL2ClientSession) GetBlockHash(number *big.Int) ([32]byte, error) {
	return _TaikoL2Client.Contract.GetBlockHash(&_TaikoL2Client.CallOpts, number)
}

// GetBlockHash is a free data retrieval call binding the contract method 0xee82ac5e.
//
// Solidity: function getBlockHash(uint256 number) view returns(bytes32)
func (_TaikoL2Client *TaikoL2ClientCallerSession) GetBlockHash(number *big.Int) ([32]byte, error) {
	return _TaikoL2Client.Contract.GetBlockHash(&_TaikoL2Client.CallOpts, number)
}

// GetXchainBlockHash is a free data retrieval call binding the contract method 0xa4e6775f.
//
// Solidity: function getXchainBlockHash(uint256 number) view returns(bytes32)
func (_TaikoL2Client *TaikoL2ClientCaller) GetXchainBlockHash(opts *bind.CallOpts, number *big.Int) ([32]byte, error) {
	var out []interface{}
	err := _TaikoL2Client.contract.Call(opts, &out, "getXchainBlockHash", number)

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// GetXchainBlockHash is a free data retrieval call binding the contract method 0xa4e6775f.
//
// Solidity: function getXchainBlockHash(uint256 number) view returns(bytes32)
func (_TaikoL2Client *TaikoL2ClientSession) GetXchainBlockHash(number *big.Int) ([32]byte, error) {
	return _TaikoL2Client.Contract.GetXchainBlockHash(&_TaikoL2Client.CallOpts, number)
}

// GetXchainBlockHash is a free data retrieval call binding the contract method 0xa4e6775f.
//
// Solidity: function getXchainBlockHash(uint256 number) view returns(bytes32)
func (_TaikoL2Client *TaikoL2ClientCallerSession) GetXchainBlockHash(number *big.Int) ([32]byte, error) {
	return _TaikoL2Client.Contract.GetXchainBlockHash(&_TaikoL2Client.CallOpts, number)
}

// GetXchainSignalRoot is a free data retrieval call binding the contract method 0x609bbd06.
//
// Solidity: function getXchainSignalRoot(uint256 number) view returns(bytes32)
func (_TaikoL2Client *TaikoL2ClientCaller) GetXchainSignalRoot(opts *bind.CallOpts, number *big.Int) ([32]byte, error) {
	var out []interface{}
	err := _TaikoL2Client.contract.Call(opts, &out, "getXchainSignalRoot", number)

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// GetXchainSignalRoot is a free data retrieval call binding the contract method 0x609bbd06.
//
// Solidity: function getXchainSignalRoot(uint256 number) view returns(bytes32)
func (_TaikoL2Client *TaikoL2ClientSession) GetXchainSignalRoot(number *big.Int) ([32]byte, error) {
	return _TaikoL2Client.Contract.GetXchainSignalRoot(&_TaikoL2Client.CallOpts, number)
}

// GetXchainSignalRoot is a free data retrieval call binding the contract method 0x609bbd06.
//
// Solidity: function getXchainSignalRoot(uint256 number) view returns(bytes32)
func (_TaikoL2Client *TaikoL2ClientCallerSession) GetXchainSignalRoot(number *big.Int) ([32]byte, error) {
	return _TaikoL2Client.Contract.GetXchainSignalRoot(&_TaikoL2Client.CallOpts, number)
}

// KeyForName is a free data retrieval call binding the contract method 0x3e98a12e.
//
// Solidity: function keyForName(uint256 chainId, string name) pure returns(string)
func (_TaikoL2Client *TaikoL2ClientCaller) KeyForName(opts *bind.CallOpts, chainId *big.Int, name string) (string, error) {
	var out []interface{}
	err := _TaikoL2Client.contract.Call(opts, &out, "keyForName", chainId, name)

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// KeyForName is a free data retrieval call binding the contract method 0x3e98a12e.
//
// Solidity: function keyForName(uint256 chainId, string name) pure returns(string)
func (_TaikoL2Client *TaikoL2ClientSession) KeyForName(chainId *big.Int, name string) (string, error) {
	return _TaikoL2Client.Contract.KeyForName(&_TaikoL2Client.CallOpts, chainId, name)
}

// KeyForName is a free data retrieval call binding the contract method 0x3e98a12e.
//
// Solidity: function keyForName(uint256 chainId, string name) pure returns(string)
func (_TaikoL2Client *TaikoL2ClientCallerSession) KeyForName(chainId *big.Int, name string) (string, error) {
	return _TaikoL2Client.Contract.KeyForName(&_TaikoL2Client.CallOpts, chainId, name)
}

// LatestSyncedL1Height is a free data retrieval call binding the contract method 0xc7b96908.
//
// Solidity: function latestSyncedL1Height() view returns(uint64)
func (_TaikoL2Client *TaikoL2ClientCaller) LatestSyncedL1Height(opts *bind.CallOpts) (uint64, error) {
	var out []interface{}
	err := _TaikoL2Client.contract.Call(opts, &out, "latestSyncedL1Height")

	if err != nil {
		return *new(uint64), err
	}

	out0 := *abi.ConvertType(out[0], new(uint64)).(*uint64)

	return out0, err

}

// LatestSyncedL1Height is a free data retrieval call binding the contract method 0xc7b96908.
//
// Solidity: function latestSyncedL1Height() view returns(uint64)
func (_TaikoL2Client *TaikoL2ClientSession) LatestSyncedL1Height() (uint64, error) {
	return _TaikoL2Client.Contract.LatestSyncedL1Height(&_TaikoL2Client.CallOpts)
}

// LatestSyncedL1Height is a free data retrieval call binding the contract method 0xc7b96908.
//
// Solidity: function latestSyncedL1Height() view returns(uint64)
func (_TaikoL2Client *TaikoL2ClientCallerSession) LatestSyncedL1Height() (uint64, error) {
	return _TaikoL2Client.Contract.LatestSyncedL1Height(&_TaikoL2Client.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_TaikoL2Client *TaikoL2ClientCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _TaikoL2Client.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_TaikoL2Client *TaikoL2ClientSession) Owner() (common.Address, error) {
	return _TaikoL2Client.Contract.Owner(&_TaikoL2Client.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_TaikoL2Client *TaikoL2ClientCallerSession) Owner() (common.Address, error) {
	return _TaikoL2Client.Contract.Owner(&_TaikoL2Client.CallOpts)
}

// ParentTimestamp is a free data retrieval call binding the contract method 0x539b8ade.
//
// Solidity: function parentTimestamp() view returns(uint64)
func (_TaikoL2Client *TaikoL2ClientCaller) ParentTimestamp(opts *bind.CallOpts) (uint64, error) {
	var out []interface{}
	err := _TaikoL2Client.contract.Call(opts, &out, "parentTimestamp")

	if err != nil {
		return *new(uint64), err
	}

	out0 := *abi.ConvertType(out[0], new(uint64)).(*uint64)

	return out0, err

}

// ParentTimestamp is a free data retrieval call binding the contract method 0x539b8ade.
//
// Solidity: function parentTimestamp() view returns(uint64)
func (_TaikoL2Client *TaikoL2ClientSession) ParentTimestamp() (uint64, error) {
	return _TaikoL2Client.Contract.ParentTimestamp(&_TaikoL2Client.CallOpts)
}

// ParentTimestamp is a free data retrieval call binding the contract method 0x539b8ade.
//
// Solidity: function parentTimestamp() view returns(uint64)
func (_TaikoL2Client *TaikoL2ClientCallerSession) ParentTimestamp() (uint64, error) {
	return _TaikoL2Client.Contract.ParentTimestamp(&_TaikoL2Client.CallOpts)
}

// PublicInputHash is a free data retrieval call binding the contract method 0xdac5df78.
//
// Solidity: function publicInputHash() view returns(bytes32)
func (_TaikoL2Client *TaikoL2ClientCaller) PublicInputHash(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _TaikoL2Client.contract.Call(opts, &out, "publicInputHash")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// PublicInputHash is a free data retrieval call binding the contract method 0xdac5df78.
//
// Solidity: function publicInputHash() view returns(bytes32)
func (_TaikoL2Client *TaikoL2ClientSession) PublicInputHash() ([32]byte, error) {
	return _TaikoL2Client.Contract.PublicInputHash(&_TaikoL2Client.CallOpts)
}

// PublicInputHash is a free data retrieval call binding the contract method 0xdac5df78.
//
// Solidity: function publicInputHash() view returns(bytes32)
func (_TaikoL2Client *TaikoL2ClientCallerSession) PublicInputHash() ([32]byte, error) {
	return _TaikoL2Client.Contract.PublicInputHash(&_TaikoL2Client.CallOpts)
}

// Resolve is a free data retrieval call binding the contract method 0x0ca4dffd.
//
// Solidity: function resolve(string name, bool allowZeroAddress) view returns(address)
func (_TaikoL2Client *TaikoL2ClientCaller) Resolve(opts *bind.CallOpts, name string, allowZeroAddress bool) (common.Address, error) {
	var out []interface{}
	err := _TaikoL2Client.contract.Call(opts, &out, "resolve", name, allowZeroAddress)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Resolve is a free data retrieval call binding the contract method 0x0ca4dffd.
//
// Solidity: function resolve(string name, bool allowZeroAddress) view returns(address)
func (_TaikoL2Client *TaikoL2ClientSession) Resolve(name string, allowZeroAddress bool) (common.Address, error) {
	return _TaikoL2Client.Contract.Resolve(&_TaikoL2Client.CallOpts, name, allowZeroAddress)
}

// Resolve is a free data retrieval call binding the contract method 0x0ca4dffd.
//
// Solidity: function resolve(string name, bool allowZeroAddress) view returns(address)
func (_TaikoL2Client *TaikoL2ClientCallerSession) Resolve(name string, allowZeroAddress bool) (common.Address, error) {
	return _TaikoL2Client.Contract.Resolve(&_TaikoL2Client.CallOpts, name, allowZeroAddress)
}

// Resolve0 is a free data retrieval call binding the contract method 0x1be2bfa7.
//
// Solidity: function resolve(uint256 chainId, string name, bool allowZeroAddress) view returns(address)
func (_TaikoL2Client *TaikoL2ClientCaller) Resolve0(opts *bind.CallOpts, chainId *big.Int, name string, allowZeroAddress bool) (common.Address, error) {
	var out []interface{}
	err := _TaikoL2Client.contract.Call(opts, &out, "resolve0", chainId, name, allowZeroAddress)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Resolve0 is a free data retrieval call binding the contract method 0x1be2bfa7.
//
// Solidity: function resolve(uint256 chainId, string name, bool allowZeroAddress) view returns(address)
func (_TaikoL2Client *TaikoL2ClientSession) Resolve0(chainId *big.Int, name string, allowZeroAddress bool) (common.Address, error) {
	return _TaikoL2Client.Contract.Resolve0(&_TaikoL2Client.CallOpts, chainId, name, allowZeroAddress)
}

// Resolve0 is a free data retrieval call binding the contract method 0x1be2bfa7.
//
// Solidity: function resolve(uint256 chainId, string name, bool allowZeroAddress) view returns(address)
func (_TaikoL2Client *TaikoL2ClientCallerSession) Resolve0(chainId *big.Int, name string, allowZeroAddress bool) (common.Address, error) {
	return _TaikoL2Client.Contract.Resolve0(&_TaikoL2Client.CallOpts, chainId, name, allowZeroAddress)
}

// SignAnchor is a free data retrieval call binding the contract method 0x591aad8a.
//
// Solidity: function signAnchor(bytes32 digest, uint8 k) view returns(uint8 v, uint256 r, uint256 s)
func (_TaikoL2Client *TaikoL2ClientCaller) SignAnchor(opts *bind.CallOpts, digest [32]byte, k uint8) (struct {
	V uint8
	R *big.Int
	S *big.Int
}, error) {
	var out []interface{}
	err := _TaikoL2Client.contract.Call(opts, &out, "signAnchor", digest, k)

	outstruct := new(struct {
		V uint8
		R *big.Int
		S *big.Int
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.V = *abi.ConvertType(out[0], new(uint8)).(*uint8)
	outstruct.R = *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)
	outstruct.S = *abi.ConvertType(out[2], new(*big.Int)).(**big.Int)

	return *outstruct, err

}

// SignAnchor is a free data retrieval call binding the contract method 0x591aad8a.
//
// Solidity: function signAnchor(bytes32 digest, uint8 k) view returns(uint8 v, uint256 r, uint256 s)
func (_TaikoL2Client *TaikoL2ClientSession) SignAnchor(digest [32]byte, k uint8) (struct {
	V uint8
	R *big.Int
	S *big.Int
}, error) {
	return _TaikoL2Client.Contract.SignAnchor(&_TaikoL2Client.CallOpts, digest, k)
}

// SignAnchor is a free data retrieval call binding the contract method 0x591aad8a.
//
// Solidity: function signAnchor(bytes32 digest, uint8 k) view returns(uint8 v, uint256 r, uint256 s)
func (_TaikoL2Client *TaikoL2ClientCallerSession) SignAnchor(digest [32]byte, k uint8) (struct {
	V uint8
	R *big.Int
	S *big.Int
}, error) {
	return _TaikoL2Client.Contract.SignAnchor(&_TaikoL2Client.CallOpts, digest, k)
}

// Xscale is a free data retrieval call binding the contract method 0xf5c97740.
//
// Solidity: function xscale() view returns(uint64)
func (_TaikoL2Client *TaikoL2ClientCaller) Xscale(opts *bind.CallOpts) (uint64, error) {
	var out []interface{}
	err := _TaikoL2Client.contract.Call(opts, &out, "xscale")

	if err != nil {
		return *new(uint64), err
	}

	out0 := *abi.ConvertType(out[0], new(uint64)).(*uint64)

	return out0, err

}

// Xscale is a free data retrieval call binding the contract method 0xf5c97740.
//
// Solidity: function xscale() view returns(uint64)
func (_TaikoL2Client *TaikoL2ClientSession) Xscale() (uint64, error) {
	return _TaikoL2Client.Contract.Xscale(&_TaikoL2Client.CallOpts)
}

// Xscale is a free data retrieval call binding the contract method 0xf5c97740.
//
// Solidity: function xscale() view returns(uint64)
func (_TaikoL2Client *TaikoL2ClientCallerSession) Xscale() (uint64, error) {
	return _TaikoL2Client.Contract.Xscale(&_TaikoL2Client.CallOpts)
}

// Yscale is a free data retrieval call binding the contract method 0x3fa85350.
//
// Solidity: function yscale() view returns(uint128)
func (_TaikoL2Client *TaikoL2ClientCaller) Yscale(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _TaikoL2Client.contract.Call(opts, &out, "yscale")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Yscale is a free data retrieval call binding the contract method 0x3fa85350.
//
// Solidity: function yscale() view returns(uint128)
func (_TaikoL2Client *TaikoL2ClientSession) Yscale() (*big.Int, error) {
	return _TaikoL2Client.Contract.Yscale(&_TaikoL2Client.CallOpts)
}

// Yscale is a free data retrieval call binding the contract method 0x3fa85350.
//
// Solidity: function yscale() view returns(uint128)
func (_TaikoL2Client *TaikoL2ClientCallerSession) Yscale() (*big.Int, error) {
	return _TaikoL2Client.Contract.Yscale(&_TaikoL2Client.CallOpts)
}

// Anchor is a paid mutator transaction binding the contract method 0x3d384a4b.
//
// Solidity: function anchor(bytes32 l1Hash, bytes32 l1SignalRoot, uint64 l1Height, uint64 parentGasUsed) returns()
func (_TaikoL2Client *TaikoL2ClientTransactor) Anchor(opts *bind.TransactOpts, l1Hash [32]byte, l1SignalRoot [32]byte, l1Height uint64, parentGasUsed uint64) (*types.Transaction, error) {
	return _TaikoL2Client.contract.Transact(opts, "anchor", l1Hash, l1SignalRoot, l1Height, parentGasUsed)
}

// Anchor is a paid mutator transaction binding the contract method 0x3d384a4b.
//
// Solidity: function anchor(bytes32 l1Hash, bytes32 l1SignalRoot, uint64 l1Height, uint64 parentGasUsed) returns()
func (_TaikoL2Client *TaikoL2ClientSession) Anchor(l1Hash [32]byte, l1SignalRoot [32]byte, l1Height uint64, parentGasUsed uint64) (*types.Transaction, error) {
	return _TaikoL2Client.Contract.Anchor(&_TaikoL2Client.TransactOpts, l1Hash, l1SignalRoot, l1Height, parentGasUsed)
}

// Anchor is a paid mutator transaction binding the contract method 0x3d384a4b.
//
// Solidity: function anchor(bytes32 l1Hash, bytes32 l1SignalRoot, uint64 l1Height, uint64 parentGasUsed) returns()
func (_TaikoL2Client *TaikoL2ClientTransactorSession) Anchor(l1Hash [32]byte, l1SignalRoot [32]byte, l1Height uint64, parentGasUsed uint64) (*types.Transaction, error) {
	return _TaikoL2Client.Contract.Anchor(&_TaikoL2Client.TransactOpts, l1Hash, l1SignalRoot, l1Height, parentGasUsed)
}

// Init is a paid mutator transaction binding the contract method 0x8f3ca30d.
//
// Solidity: function init(address _addressManager, (uint64,uint64,uint64,uint64,uint64) _param1559) returns()
func (_TaikoL2Client *TaikoL2ClientTransactor) Init(opts *bind.TransactOpts, _addressManager common.Address, _param1559 TaikoL2EIP1559Params) (*types.Transaction, error) {
	return _TaikoL2Client.contract.Transact(opts, "init", _addressManager, _param1559)
}

// Init is a paid mutator transaction binding the contract method 0x8f3ca30d.
//
// Solidity: function init(address _addressManager, (uint64,uint64,uint64,uint64,uint64) _param1559) returns()
func (_TaikoL2Client *TaikoL2ClientSession) Init(_addressManager common.Address, _param1559 TaikoL2EIP1559Params) (*types.Transaction, error) {
	return _TaikoL2Client.Contract.Init(&_TaikoL2Client.TransactOpts, _addressManager, _param1559)
}

// Init is a paid mutator transaction binding the contract method 0x8f3ca30d.
//
// Solidity: function init(address _addressManager, (uint64,uint64,uint64,uint64,uint64) _param1559) returns()
func (_TaikoL2Client *TaikoL2ClientTransactorSession) Init(_addressManager common.Address, _param1559 TaikoL2EIP1559Params) (*types.Transaction, error) {
	return _TaikoL2Client.Contract.Init(&_TaikoL2Client.TransactOpts, _addressManager, _param1559)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_TaikoL2Client *TaikoL2ClientTransactor) RenounceOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _TaikoL2Client.contract.Transact(opts, "renounceOwnership")
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_TaikoL2Client *TaikoL2ClientSession) RenounceOwnership() (*types.Transaction, error) {
	return _TaikoL2Client.Contract.RenounceOwnership(&_TaikoL2Client.TransactOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_TaikoL2Client *TaikoL2ClientTransactorSession) RenounceOwnership() (*types.Transaction, error) {
	return _TaikoL2Client.Contract.RenounceOwnership(&_TaikoL2Client.TransactOpts)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_TaikoL2Client *TaikoL2ClientTransactor) TransferOwnership(opts *bind.TransactOpts, newOwner common.Address) (*types.Transaction, error) {
	return _TaikoL2Client.contract.Transact(opts, "transferOwnership", newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_TaikoL2Client *TaikoL2ClientSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _TaikoL2Client.Contract.TransferOwnership(&_TaikoL2Client.TransactOpts, newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_TaikoL2Client *TaikoL2ClientTransactorSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _TaikoL2Client.Contract.TransferOwnership(&_TaikoL2Client.TransactOpts, newOwner)
}

// TaikoL2ClientBlockVarsIterator is returned from FilterBlockVars and is used to iterate over the raw logs and unpacked data for BlockVars events raised by the TaikoL2Client contract.
type TaikoL2ClientBlockVarsIterator struct {
	Event *TaikoL2ClientBlockVars // Event containing the contract specifics and raw log

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
func (it *TaikoL2ClientBlockVarsIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TaikoL2ClientBlockVars)
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
		it.Event = new(TaikoL2ClientBlockVars)
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
func (it *TaikoL2ClientBlockVarsIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *TaikoL2ClientBlockVarsIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// TaikoL2ClientBlockVars represents a BlockVars event raised by the TaikoL2Client contract.
type TaikoL2ClientBlockVars struct {
	Number     uint64
	Basefee    uint64
	Gaslimit   uint64
	Timestamp  uint64
	ParentHash [32]byte
	Prevrandao *big.Int
	Coinbase   common.Address
	Chainid    uint32
	Raw        types.Log // Blockchain specific contextual infos
}

// FilterBlockVars is a free log retrieval operation binding the contract event 0xc79d10e4d7bf09ff9c9e62072d505d73f62dfea7216692b2a6cfbb72c9aaa3ca.
//
// Solidity: event BlockVars(uint64 number, uint64 basefee, uint64 gaslimit, uint64 timestamp, bytes32 parentHash, uint256 prevrandao, address coinbase, uint32 chainid)
func (_TaikoL2Client *TaikoL2ClientFilterer) FilterBlockVars(opts *bind.FilterOpts) (*TaikoL2ClientBlockVarsIterator, error) {

	logs, sub, err := _TaikoL2Client.contract.FilterLogs(opts, "BlockVars")
	if err != nil {
		return nil, err
	}
	return &TaikoL2ClientBlockVarsIterator{contract: _TaikoL2Client.contract, event: "BlockVars", logs: logs, sub: sub}, nil
}

// WatchBlockVars is a free log subscription operation binding the contract event 0xc79d10e4d7bf09ff9c9e62072d505d73f62dfea7216692b2a6cfbb72c9aaa3ca.
//
// Solidity: event BlockVars(uint64 number, uint64 basefee, uint64 gaslimit, uint64 timestamp, bytes32 parentHash, uint256 prevrandao, address coinbase, uint32 chainid)
func (_TaikoL2Client *TaikoL2ClientFilterer) WatchBlockVars(opts *bind.WatchOpts, sink chan<- *TaikoL2ClientBlockVars) (event.Subscription, error) {

	logs, sub, err := _TaikoL2Client.contract.WatchLogs(opts, "BlockVars")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(TaikoL2ClientBlockVars)
				if err := _TaikoL2Client.contract.UnpackLog(event, "BlockVars", log); err != nil {
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

// ParseBlockVars is a log parse operation binding the contract event 0xc79d10e4d7bf09ff9c9e62072d505d73f62dfea7216692b2a6cfbb72c9aaa3ca.
//
// Solidity: event BlockVars(uint64 number, uint64 basefee, uint64 gaslimit, uint64 timestamp, bytes32 parentHash, uint256 prevrandao, address coinbase, uint32 chainid)
func (_TaikoL2Client *TaikoL2ClientFilterer) ParseBlockVars(log types.Log) (*TaikoL2ClientBlockVars, error) {
	event := new(TaikoL2ClientBlockVars)
	if err := _TaikoL2Client.contract.UnpackLog(event, "BlockVars", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// TaikoL2ClientInitializedIterator is returned from FilterInitialized and is used to iterate over the raw logs and unpacked data for Initialized events raised by the TaikoL2Client contract.
type TaikoL2ClientInitializedIterator struct {
	Event *TaikoL2ClientInitialized // Event containing the contract specifics and raw log

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
func (it *TaikoL2ClientInitializedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TaikoL2ClientInitialized)
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
		it.Event = new(TaikoL2ClientInitialized)
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
func (it *TaikoL2ClientInitializedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *TaikoL2ClientInitializedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// TaikoL2ClientInitialized represents a Initialized event raised by the TaikoL2Client contract.
type TaikoL2ClientInitialized struct {
	Version uint8
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterInitialized is a free log retrieval operation binding the contract event 0x7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb3847402498.
//
// Solidity: event Initialized(uint8 version)
func (_TaikoL2Client *TaikoL2ClientFilterer) FilterInitialized(opts *bind.FilterOpts) (*TaikoL2ClientInitializedIterator, error) {

	logs, sub, err := _TaikoL2Client.contract.FilterLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return &TaikoL2ClientInitializedIterator{contract: _TaikoL2Client.contract, event: "Initialized", logs: logs, sub: sub}, nil
}

// WatchInitialized is a free log subscription operation binding the contract event 0x7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb3847402498.
//
// Solidity: event Initialized(uint8 version)
func (_TaikoL2Client *TaikoL2ClientFilterer) WatchInitialized(opts *bind.WatchOpts, sink chan<- *TaikoL2ClientInitialized) (event.Subscription, error) {

	logs, sub, err := _TaikoL2Client.contract.WatchLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(TaikoL2ClientInitialized)
				if err := _TaikoL2Client.contract.UnpackLog(event, "Initialized", log); err != nil {
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

// ParseInitialized is a log parse operation binding the contract event 0x7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb3847402498.
//
// Solidity: event Initialized(uint8 version)
func (_TaikoL2Client *TaikoL2ClientFilterer) ParseInitialized(log types.Log) (*TaikoL2ClientInitialized, error) {
	event := new(TaikoL2ClientInitialized)
	if err := _TaikoL2Client.contract.UnpackLog(event, "Initialized", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// TaikoL2ClientOwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the TaikoL2Client contract.
type TaikoL2ClientOwnershipTransferredIterator struct {
	Event *TaikoL2ClientOwnershipTransferred // Event containing the contract specifics and raw log

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
func (it *TaikoL2ClientOwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TaikoL2ClientOwnershipTransferred)
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
		it.Event = new(TaikoL2ClientOwnershipTransferred)
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
func (it *TaikoL2ClientOwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *TaikoL2ClientOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// TaikoL2ClientOwnershipTransferred represents a OwnershipTransferred event raised by the TaikoL2Client contract.
type TaikoL2ClientOwnershipTransferred struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_TaikoL2Client *TaikoL2ClientFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*TaikoL2ClientOwnershipTransferredIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _TaikoL2Client.contract.FilterLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &TaikoL2ClientOwnershipTransferredIterator{contract: _TaikoL2Client.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_TaikoL2Client *TaikoL2ClientFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *TaikoL2ClientOwnershipTransferred, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _TaikoL2Client.contract.WatchLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(TaikoL2ClientOwnershipTransferred)
				if err := _TaikoL2Client.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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

// ParseOwnershipTransferred is a log parse operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_TaikoL2Client *TaikoL2ClientFilterer) ParseOwnershipTransferred(log types.Log) (*TaikoL2ClientOwnershipTransferred, error) {
	event := new(TaikoL2ClientOwnershipTransferred)
	if err := _TaikoL2Client.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// TaikoL2ClientXchainSyncedIterator is returned from FilterXchainSynced and is used to iterate over the raw logs and unpacked data for XchainSynced events raised by the TaikoL2Client contract.
type TaikoL2ClientXchainSyncedIterator struct {
	Event *TaikoL2ClientXchainSynced // Event containing the contract specifics and raw log

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
func (it *TaikoL2ClientXchainSyncedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TaikoL2ClientXchainSynced)
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
		it.Event = new(TaikoL2ClientXchainSynced)
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
func (it *TaikoL2ClientXchainSyncedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *TaikoL2ClientXchainSyncedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// TaikoL2ClientXchainSynced represents a XchainSynced event raised by the TaikoL2Client contract.
type TaikoL2ClientXchainSynced struct {
	SrcHeight  *big.Int
	BlockHash  [32]byte
	SignalRoot [32]byte
	Raw        types.Log // Blockchain specific contextual infos
}

// FilterXchainSynced is a free log retrieval operation binding the contract event 0xc7edd3d480c294297f3924d0ffab64074e7fb22e004ea492d5dd691fa1fc99c0.
//
// Solidity: event XchainSynced(uint256 indexed srcHeight, bytes32 blockHash, bytes32 signalRoot)
func (_TaikoL2Client *TaikoL2ClientFilterer) FilterXchainSynced(opts *bind.FilterOpts, srcHeight []*big.Int) (*TaikoL2ClientXchainSyncedIterator, error) {

	var srcHeightRule []interface{}
	for _, srcHeightItem := range srcHeight {
		srcHeightRule = append(srcHeightRule, srcHeightItem)
	}

	logs, sub, err := _TaikoL2Client.contract.FilterLogs(opts, "XchainSynced", srcHeightRule)
	if err != nil {
		return nil, err
	}
	return &TaikoL2ClientXchainSyncedIterator{contract: _TaikoL2Client.contract, event: "XchainSynced", logs: logs, sub: sub}, nil
}

// WatchXchainSynced is a free log subscription operation binding the contract event 0xc7edd3d480c294297f3924d0ffab64074e7fb22e004ea492d5dd691fa1fc99c0.
//
// Solidity: event XchainSynced(uint256 indexed srcHeight, bytes32 blockHash, bytes32 signalRoot)
func (_TaikoL2Client *TaikoL2ClientFilterer) WatchXchainSynced(opts *bind.WatchOpts, sink chan<- *TaikoL2ClientXchainSynced, srcHeight []*big.Int) (event.Subscription, error) {

	var srcHeightRule []interface{}
	for _, srcHeightItem := range srcHeight {
		srcHeightRule = append(srcHeightRule, srcHeightItem)
	}

	logs, sub, err := _TaikoL2Client.contract.WatchLogs(opts, "XchainSynced", srcHeightRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(TaikoL2ClientXchainSynced)
				if err := _TaikoL2Client.contract.UnpackLog(event, "XchainSynced", log); err != nil {
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

// ParseXchainSynced is a log parse operation binding the contract event 0xc7edd3d480c294297f3924d0ffab64074e7fb22e004ea492d5dd691fa1fc99c0.
//
// Solidity: event XchainSynced(uint256 indexed srcHeight, bytes32 blockHash, bytes32 signalRoot)
func (_TaikoL2Client *TaikoL2ClientFilterer) ParseXchainSynced(log types.Log) (*TaikoL2ClientXchainSynced, error) {
	event := new(TaikoL2ClientXchainSynced)
	if err := _TaikoL2Client.contract.UnpackLog(event, "XchainSynced", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
