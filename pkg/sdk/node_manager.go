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
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	nmabi "github.com/ethereum/go-ethereum/contracts/native/go_abi/node_manager_abi"
	nm "github.com/ethereum/go-ethereum/contracts/native/governance/node_manager"
	"github.com/ethereum/go-ethereum/contracts/native/utils"
	"github.com/ethereum/go-ethereum/rlp"
)

var (
	nodeManagerAddr = utils.NodeManagerContractAddress
)

func init() {
	nm.InitABI()
}

func (c *Account) Epoch() (*nm.EpochInfo, error) {
	payload, err := new(nm.GetCurrentEpochInfoParam).Encode()
	if err != nil {
		return nil, err
	}
	output, err := c.callNodeManager(payload, nil)
	if err != nil {
		return nil, err
	}

	var (
		raw   []byte
		epoch = new(nm.EpochInfo)
	)
	if err := utils.UnpackOutputs(nm.ABI, nmabi.MethodGetCurrentEpochInfo, &raw, output); err != nil {
		return nil, err
	}
	if err := rlp.DecodeBytes(raw, epoch); err != nil {
		return nil, err
	}
	return epoch, nil
}

func (c *Account) Register(validator common.Address, amount *big.Int, desc string) (common.Hash, error) {
	input := &nm.CreateValidatorParam{
		ConsensusAddress: validator,
		SignerAddress:    validator,
		ProposalAddress:  c.addr,
		Commission:       big.NewInt(0),
		InitStake:        amount,
		Desc:             desc,
	}

	payload, err := input.Encode()
	if err != nil {
		return common.EmptyHash, err
	}

	return c.sendNodeManagerTx(payload)
}

func (c *Account) Stake(validator common.Address, amount *big.Int) (common.Hash, error) {
	input := &nm.StakeParam{
		ConsensusAddress: validator,
		Amount:           amount,
	}
	payload, err := input.Encode()
	if err != nil {
		return common.EmptyHash, err
	}
	return c.sendNodeManagerTx(payload)
}

func (c *Account) sendNodeManagerTx(payload []byte) (common.Hash, error) {
	return c.signAndSendTx(payload, nodeManagerAddr)
}

func (c *Account) callNodeManager(payload []byte, blockNum *big.Int) ([]byte, error) {
	return c.CallContract(c.Addr(), nodeManagerAddr, payload, blockNum)
}
