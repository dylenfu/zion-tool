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

package sdk

import (
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/contracts/native"
	scom "github.com/ethereum/go-ethereum/contracts/native/cross_chain_manager/common"
	mlp "github.com/ethereum/go-ethereum/contracts/native/cross_chain_manager/zion/mainchain/lock_proxy"
	slp "github.com/ethereum/go-ethereum/contracts/native/cross_chain_manager/zion/sidechain/lock_proxy"
	"github.com/ethereum/go-ethereum/contracts/native/utils"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/rlp"
)

func init() {
	mlp.InitABI()
	slp.InitABI()
}

func (c *Account) Mint(targetCrossChainID uint64, to common.Address, amount *big.Int) (common.Hash, error) {
	input := &mlp.MethodLockInput{
		ToChainId: targetCrossChainID,
		ToAddress: to,
		Amount:    amount,
	}
	payload, err := input.Encode()
	if err != nil {
		return common.EmptyHash, err
	}

	return c.sendLockProxyTx(payload, amount)
}

func (c *Account) Burn(amount *big.Int) (common.Hash, error) {
	input := &slp.MethodBurnInput{
		ToChainId: native.ZionMainChainID,
		Amount:    amount,
	}

	payload, err := input.Encode()
	if err != nil {
		return common.EmptyHash, err
	}

	return c.sendLockProxyTx(payload, amount)
}

//state := &nccmc.EntranceParam{
//	SourceChainID:         sourceChainId,
//	Height:                height,
//	Proof:                 proof,
//	RelayerAddress:        relayerAddress,
//	Extra:                 txData,
//	HeaderOrCrossChainMsg: HeaderOrCrossChainMsg,
//}
func (c *Account) ImportOutTransfer(sourceChainID uint64, height uint32, txData, proof, rawHeader []byte) (common.Hash, error) {
	relayerAddr := c.address[:]
	payload, err := utils.PackMethod(scom.ABI, scom.MethodImportOuterTransfer, sourceChainID, height, proof, relayerAddr, txData, rawHeader)
	if err != nil {
		return common.EmptyHash, err
	}

	return c.signAndSendTx(payload, utils.CrossChainManagerContractAddress)
}

func (c *Account) sendLockProxyTx(payload []byte, amount *big.Int) (common.Hash, error) {
	return c.signAndSendTxWithValue(payload, amount, utils.LockProxyContractAddress)
}

func (c *Account) callLockProxy(payload []byte, blockNum string) ([]byte, error) {
	return c.CallContract(c.Address(), utils.LockProxyContractAddress, payload, blockNum)
}

func (c *Account) GetRawHeaderAndSeals(number uint64) (*types.Header, []byte, []byte, error) {
	header, err := c.BlockHeaderByNumber(number)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to fetch header, err: %v", err)
	}
	rawHeader, err := rlp.EncodeToBytes(types.HotstuffFilteredHeader(header, false))
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to encode header, err: %v", err)
	}
	extra, err := types.ExtractHotstuffExtra(header)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to extra header, err: %v", err)
	}
	rawSeals, err := rlp.EncodeToBytes(extra.CommittedSeal)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to encode committed seals, err: %v", err)
	}
	return header, rawHeader, rawSeals, nil
}

/*
function verifyHeaderAndExecuteTx(
        bytes memory rawHeader,
        bytes memory rawSeals,
        bytes memory accountProof,
        bytes memory storageProof,
        bytes memory rawCrossTx
    ) public returns (bool)

// event CrossChainEvent(
	address indexed sender,
	bytes txId,
	address proxyOrAssetContract,
	uint64 toChainId,
	bytes toContract,
	bytes rawdata);
*/

var (
	BytesTy, _   = abi.NewType("bytes", "", nil)
	AddrTy, _    = abi.NewType("address", "", nil)
	Uint64Ty, _  = abi.NewType("uint64", "", nil)
	Uint256Ty, _ = abi.NewType("uint256", "", nil)
)

const (
	methodVerifyHeaderAndExecuteTx = "verifyHeaderAndExecuteTx"
	eventCrossChain                = "CrossChainEvent"
)

var (
	midVerifyHeaderAndExecuteTx = crypto.Keccak256(utils.EncodePacked([]byte(methodVerifyHeaderAndExecuteTx), []byte("(bytes,bytes,bytes,bytes,bytes)")))[:4]
	eidCrossChain               = crypto.Keccak256(utils.EncodePacked([]byte(eventCrossChain), []byte("(address, bytes, address, uint64, bytes, bytes)")))[:4]

	argsVerifyHeaderAndExecuteTx = abi.Arguments{
		{Type: BytesTy, Name: "rawHeader"},
		{Type: BytesTy, Name: "rawSeals"},
		{Type: BytesTy, Name: "accountProof"},
		{Type: BytesTy, Name: "storageProof"},
		{Type: BytesTy, Name: "rawCrossTx"},
	}

	argsCrossChain = abi.Arguments{
		{Type: AddrTy, Name: "sender", Indexed: true},
		{Type: BytesTy, Name: "txId"},
		{Type: AddrTy, Name: "proxyOrAssetContract"},
		{Type: Uint64Ty, Name: "toChainId"},
		{Type: BytesTy, Name: "toContract"},
		{Type: BytesTy, Name: "rawdata"},
	}
)

func (c *Account) SideChainVerifyHeaderAndExecute(eccm common.Address, rawHeader, rawSeals, accountProof, storageProof, rawCrossTx []byte) (common.Hash, error) {
	callData, err := argsVerifyHeaderAndExecuteTx.Pack(rawHeader, rawSeals, accountProof, storageProof, rawCrossTx)
	if err != nil {
		return common.EmptyHash, err
	}
	payload := utils.EncodePacked(midVerifyHeaderAndExecuteTx, callData)

	return c.signAndSendTx(payload, eccm)
}

func UnpackSideChainCrossChainEvent(receipt *types.Log) (
	sender common.Address, txId []byte, proxyOrAsset common.Address, toChainID uint64, toContract []byte, rawData []byte, err error) {

	if len(receipt.Topics) != 2 {
		err = fmt.Errorf("event log contains indexed field, topics length should be 2")
		return
	}
	if receipt.Data == nil || len(receipt.Data) < 4 {
		err = fmt.Errorf("recepit data is invalid")
		return
	}

	// indexed field will be filtered in args.Unpack
	sender = common.BytesToAddress(receipt.Topics[1][:])

	var (
		list []interface{}
		ok   bool
	)
	if list, err = argsCrossChain.Unpack(receipt.Data); err != nil {
		return
	}
	if len(list) != 5 {
		err = fmt.Errorf("list length should be 4")
		return
	}

	if txId, ok = list[0].([]byte); !ok {
		err = fmt.Errorf("the 1st item should be bytes")
		return
	}
	if proxyOrAsset, ok = list[1].(common.Address); !ok {
		err = fmt.Errorf("the 2nd item should be address")
		return
	}
	if toChainID, ok = list[2].(uint64); !ok {
		err = fmt.Errorf("the 3rd item should be uint64")
		return
	}
	if toContract, ok = list[3].([]byte); !ok {
		err = fmt.Errorf("the 4th item should be bytes")
		return
	}
	if rawData, ok = list[4].([]byte); !ok {
		err = fmt.Errorf("the 5th item should be bytes")
		return
	}

	return
}
