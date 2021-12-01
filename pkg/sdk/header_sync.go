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
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum/common"

	"github.com/ethereum/go-ethereum"
	hsc "github.com/ethereum/go-ethereum/contracts/native/header_sync/common"
	"github.com/ethereum/go-ethereum/contracts/native/utils"
)

func init() {
	hsc.ABI = hsc.GetABI()
}

func (c *Account) SyncGenesisHeader(chainID uint64, raw []byte) (common.Hash, error) {
	payload, err := utils.PackMethod(hsc.ABI, hsc.MethodSyncGenesisHeader, chainID, raw)
	if err != nil {
		return common.EmptyHash, err
	}
	return c.sendHeaderSyncManagerTx(payload)
}

func (c *Account) SyncBlockHeader(chainID uint64, raw [][]byte) (common.Hash, error) {
	payload, err := utils.PackMethod(hsc.ABI, hsc.MethodSyncBlockHeader, chainID, raw)
	if err != nil {
		return common.EmptyHash, err
	}
	return c.sendHeaderSyncManagerTx(payload)
}

func (c *Account) sendHeaderSyncManagerTx(payload []byte) (common.Hash, error) {
	return c.signAndSendTx(payload, utils.HeaderSyncContractAddress)
}

func (c *Account) callHeaderSyncManager(payload []byte, blockNum string) ([]byte, error) {
	return c.CallContract(c.Address(), utils.HeaderSyncContractAddress, payload, blockNum)
}

func (c *Account) EstimateSubmitHeaders(chainId uint64, headers [][]byte) (gas uint64, err error) {
	data, err := hsc.ABI.Pack("syncBlockHeader", chainId, c.address, headers)
	if err != nil {
		return
	}

	msg := ethereum.CallMsg{From: c.address, To: &utils.HeaderSyncContractAddress, GasPrice: gasPrice, Value: big.NewInt(0), Data: data}
	gas, err = c.client.EstimateGas(context.Background(), msg)
	return
}
