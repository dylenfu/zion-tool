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
	"github.com/ethereum/go-ethereum/common"
	sm "github.com/ethereum/go-ethereum/contracts/native/governance/side_chain_manager"
	"github.com/ethereum/go-ethereum/contracts/native/utils"
	"math/big"
)

var (
	sideChainManagerAddr = utils.SideChainManagerContractAddress
)

func init() {
	sm.ABI = sm.GetABI()
}

func (c *Account) RegisterSideChain(
	sideChainName string,
	sideChainRouter uint64,
	chainID uint64,
	eccdAddr common.Address,
) error {
	if payload, err := utils.PackMethod(sm.ABI, sm.MethodRegisterSideChain,
		c.Address(),
		chainID,
		sideChainRouter,
		sideChainName,
		uint64(0),
		eccdAddr[:],
		[]byte{},
	); err != nil {
		return err
	} else {
		_, err = c.sendSideChainManagerTx(payload)
		return err
	}
}

func (c *Account) ApproveRegSideChain(chainID uint64) error {
	if payload, err := utils.PackMethod(sm.ABI, sm.MethodApproveRegisterSideChain, chainID, c.Address()); err != nil {
		return err
	} else {
		_, err = c.sendSideChainManagerTx(payload)
		return err
	}
}

func (c *Account) sendSideChainManagerTx(payload []byte) (common.Hash, error) {
	return c.signAndSendTx(payload, sideChainManagerAddr)
}

func (c *Account) callSideChainManager(payload []byte, blockNum *big.Int) ([]byte, error) {
	return c.CallContract(c.Address(), sideChainManagerAddr, payload, blockNum)
}
