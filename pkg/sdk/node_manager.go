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
	"bytes"
	"context"
	"encoding/hex"
	"fmt"
	"math/big"
	"strings"
	"time"

	"github.com/dylenfu/zion-tool/pkg/log"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	nm "github.com/ethereum/go-ethereum/contracts/native/governance/node_manager"
	"github.com/ethereum/go-ethereum/contracts/native/utils"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rlp"
)

var (
	nodeManagerAddr = utils.NodeManagerContractAddress
)

func init() {
	nm.InitABI()
}

func (c *Account) Epoch() (*nm.EpochInfo, error) {
	payload, err := new(nm.MethodEpochInput).Encode()
	if err != nil {
		return nil, err
	}
	enc, err := c.callNodeManager(payload, nil)
	if err != nil {
		return nil, err
	}
	if bytes.Equal(enc, utils.ByteFailed) {
		return nil, fmt.Errorf("call epoch failed")
	}

	output := new(nm.MethodEpochOutput)
	if err := output.Decode(enc); err != nil {
		return nil, err
	}
	return output.Epoch, nil
}

func (c *Account) Propose(startHeight uint64, peers *nm.Peers) (common.Hash, error) {
	input := &nm.MethodProposeInput{
		StartHeight: startHeight,
		Peers:       peers,
	}
	payload, err := input.Encode()
	if err != nil {
		return common.EmptyHash, err
	}
	return c.sendNodeManagerTx(payload)
}

func (c *Account) Vote(epochID uint64, proposal common.Hash) (common.Hash, error) {
	input := &nm.MethodVoteInput{
		EpochID:   epochID,
		EpochHash: proposal,
	}
	payload, err := input.Encode()
	if err != nil {
		return common.EmptyHash, err
	}

	return c.sendNodeManagerTx(payload)
}

func (c *Account) GetCurrentEpoch(blockNum *big.Int) (*nm.EpochInfo, error) {
	input := new(nm.MethodEpochInput)
	payload, err := input.Encode()
	if err != nil {
		return nil, err
	}

	enc, err := c.callNodeManager(payload, blockNum)
	if err != nil {
		return nil, err
	}

	output := new(nm.MethodEpochOutput)
	if err := output.Decode(enc); err != nil {
		return nil, err
	}
	return output.Epoch, nil
}

func (c *Account) GetEpochByID(id uint64, blockNum *big.Int) (*nm.EpochInfo, error) {
	input := new(nm.MethodGetEpochByIDInput)
	input.EpochID = id

	payload, err := input.Encode()
	if err != nil {
		return nil, err
	}

	enc, err := c.callNodeManager(payload, blockNum)
	if err != nil {
		return nil, err
	}

	output := new(nm.MethodEpochOutput)
	if err := output.Decode(enc); err != nil {
		return nil, err
	}

	return output.Epoch, nil
}

func (c *Account) GetProofByID(id uint64, blockNum *big.Int) (common.Hash, error) {
	input := new(nm.MethodProofInput)
	input.EpochID = id

	payload, err := input.Encode()
	if err != nil {
		return common.EmptyHash, err
	}

	enc, err := c.callNodeManager(payload, nil)
	if err != nil {
		return common.EmptyHash, err
	}

	output := new(nm.MethodProofOutput)
	if err := output.Decode(enc); err != nil {
		return common.EmptyHash, err
	}

	return output.Hash, nil
}

func (c *Account) GetChangingEpoch(blockNum *big.Int) (*nm.EpochInfo, error) {
	input := new(nm.MethodGetChangingEpochInput)

	payload, err := input.Encode()
	if err != nil {
		return nil, err
	}

	enc, err := c.callNodeManager(payload, blockNum)
	if err != nil {
		return nil, err
	}

	output := new(nm.MethodEpochOutput)
	if err := output.Decode(enc); err != nil {
		return nil, err
	}
	return output.Epoch, nil
}

func (c *Account) sendNodeManagerTx(payload []byte) (common.Hash, error) {
	return c.signAndSendTx(payload, nodeManagerAddr)
}

func (c *Account) callNodeManager(payload []byte, blockNum *big.Int) ([]byte, error) {
	return c.CallContract(c.Address(), nodeManagerAddr, payload, blockNum)
}

func (c *Account) SendTransaction(contractAddr common.Address, payload []byte) (common.Hash, error) {
	addr := c.Address()

	nonce := c.GetNonce(addr.Hex())
	if c.nonce < nonce {
		c.nonce = nonce
	}
	log.Debugf("%s current nonce %d, valid nonce %d", addr.Hex(), c.nonce, nonce)
	tx := types.NewTransaction(
		c.nonce,
		contractAddr,
		big.NewInt(0),
		gasLimit,
		big.NewInt(2000000000),
		payload,
	)
	hash := tx.Hash()

	signedTx, err := c.SignTransaction(tx)
	if err != nil {
		return hash, err
	}
	c.nonce += 1
	return c.SendRawTransaction(hash, signedTx)
}

func (c *Account) SignTransaction(tx *types.Transaction) (string, error) {

	signer := types.EIP155Signer{}
	signedTx, err := types.SignTx(
		tx,
		signer,
		c.pk,
	)
	if err != nil {
		return "", fmt.Errorf("failed to sign tx: [%v]", err)
	}

	bz, err := rlp.EncodeToBytes(signedTx)
	if err != nil {
		return "", fmt.Errorf("failed to rlp encode bytes: [%v]", err)
	}
	return "0x" + hex.EncodeToString(bz), nil
}

func (c *Account) SendRawTransaction(hash common.Hash, signedTx string) (common.Hash, error) {
	var result common.Hash
	if err := c.rpcClient.Call(&result, "eth_sendRawTransaction", signedTx); err != nil {
		return hash, fmt.Errorf("failed to send raw transaction: [%v]", err)
	}

	return result, nil
}

func (c *Account) SendTransactionAndDumpEvent(contract common.Address, payload []byte) error {
	hash, err := c.SendTransaction(contract, payload)
	if err != nil {
		return err
	}
	time.Sleep(2)
	return c.DumpEventLog(hash)
}

func (c *Account) WaitTransaction(hash common.Hash) error {
	for {
		time.Sleep(time.Second * 1)
		_, ispending, err := c.client.TransactionByHash(context.Background(), hash)
		if err != nil {
			log.Errorf("failed to call TransactionByHash: %v", err)
			continue
		}
		if ispending == true {
			continue
		}

		if err := c.DumpEventLog(hash); err != nil {
			return err
		}
		break
	}
	return nil
}

func (c *Account) GetNonce(address string) uint64 {
	var raw string

	if err := c.rpcClient.Call(
		&raw,
		"eth_getTransactionCount",
		address,
		"latest",
	); err != nil {
		panic(fmt.Errorf("failed to get nonce: [%v]", err))
	}

	without0xStr := strings.Replace(raw, "0x", "", -1)
	bigNonce, _ := new(big.Int).SetString(without0xStr, 16)
	return bigNonce.Uint64()
}

func (c *Account) DumpEventLog(hash common.Hash) error {
	raw, err := c.GetReceipt(hash)
	if err != nil {
		return fmt.Errorf("faild to get receipt %s", hash.Hex())
	}

	if raw.Status == 0 {
		return fmt.Errorf("receipt failed %s", hash.Hex())
	}

	log.Infof("txhash %s, block height %d", hash.Hex(), raw.BlockNumber.Uint64())
	for _, event := range raw.Logs {
		log.Infof("eventlog address %s", event.Address.Hex())
		log.Infof("eventlog data %s", hexutil.Encode(event.Data))
		for i, topic := range event.Topics {
			log.Infof("eventlog topic[%d] %s", i, topic.String())
		}
	}
	return nil
}

func (c *Account) GetReceipt(hash common.Hash) (*types.Receipt, error) {
	raw := &types.Receipt{}
	if err := c.rpcClient.Call(raw, "eth_getTransactionReceipt", hash.Hex()); err != nil {
		return nil, err
	}
	return raw, nil
}
