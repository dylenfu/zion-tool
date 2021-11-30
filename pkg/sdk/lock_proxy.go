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
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rlp"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/contracts/native"
	mlp "github.com/ethereum/go-ethereum/contracts/native/cross_chain_manager/zion/mainchain/lock_proxy"
	slp "github.com/ethereum/go-ethereum/contracts/native/cross_chain_manager/zion/sidechain/lock_proxy"
	"github.com/ethereum/go-ethereum/contracts/native/utils"
)

func init() {
	mlp.InitABI()
	slp.InitABI()
}

func (c *Account) Mint(targetCrossChainID uint64, to common.Address, amount *big.Int) (common.Hash, error) {
	input := &mlp.MethodLockInput{
		FromAssetHash: common.EmptyAddress,
		ToChainId:     targetCrossChainID,
		ToAddress:     to[:],
		Amount:        amount,
	}
	payload, err := input.Encode()
	if err != nil {
		return common.EmptyHash, err
	}

	return c.sendLockProxyTx(payload, amount)
}

func (c *Account) Burn(to common.Address, amount *big.Int) (common.Hash, error) {
	input := &slp.MethodBurnInput{
		ToChainId: native.ZionMainChainID,
		ToAddress: to,
		Amount:    amount,
	}

	payload, err := input.Encode()
	if err != nil {
		return common.EmptyHash, err
	}

	return c.sendLockProxyTx(payload, amount)
}

func (c *Account) sendLockProxyTx(payload []byte, amount *big.Int) (common.Hash, error) {
	return c.sendNativeTxWithValue(payload, amount, utils.LockProxyContractAddress)
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