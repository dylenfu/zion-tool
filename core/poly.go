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

package core

import (
	"math/big"

	"github.com/dylenfu/zion-tool/config"
	"github.com/dylenfu/zion-tool/pkg/log"
	"github.com/ethereum/go-ethereum/common"
)

func RegisterSideChain() bool {
	var param struct {
		SideChainID uint64
		Name        string
		Router      uint64
		ECCD        string
	}

	if err := config.LoadParams("test_reg_side_chain.json", &param); err != nil {
		log.Errorf("failed to load params, err: %v", err)
		return false
	}

	proposer, err := masterAccount()
	if err != nil {
		log.Errorf("failed to generate proposer, err: %v", err)
		return false
	}

	if err := proposer.RegisterSideChain(param.Name, param.Router, param.SideChainID, common.HexToAddress(param.ECCD)); err != nil {
		log.Errorf("failed to register side chain, err: %v", err)
		return false
	}

	return true
}

func ApproveSideChain() bool {
	var param struct {
		SideChainID uint64
		AccountList []int
	}

	log.Info("start to approve side chain...")

	if err := config.LoadParams("test_approve_side_chain.json", &param); err != nil {
		log.Errorf("failed to load params, err: %v", err)
		return false
	}

	accs, err := generateAccounts(param.AccountList)
	if err != nil {
		log.Errorf("failed to load accounts, err: %v", err)
		return false
	}

	for _, proposer := range accs {
		if err := proposer.ApproveRegSideChain(param.SideChainID); err != nil {
			log.Errorf("failed to register side chain, err: %v", err)
			return false
		}
	}

	return true
}

func SyncGenesisHeader() bool {
	var param struct {
		SideChainID   uint64
		CrossChainID  uint64
		SideChainUrl  string
		NodeIndexList []int
	}

	log.Info("start to sync genesis header...")

	if err := config.LoadParams("test_sync_genesis_header.json", &param); err != nil {
		log.Errorf("failed to load params, err: %v", err)
		return false
	}

	sideChainSdk, err := customGenerateAccount(param.SideChainUrl, param.SideChainID, "")
	if err != nil {
		log.Errorf("failed to generate side chain sdk, err: %v", err)
		return false
	}

	genesisHeader, err := sideChainSdk.BlockHeaderByNumber(0)
	if err != nil {
		log.Errorf("failed to fetch side chain genesis header, err: %v", err)
		return false
	}
	rawHeader, err := genesisHeader.MarshalJSON()
	if err != nil {
		log.Errorf("failed to marshal side chain genesis header, err: %v", err)
		return false
	}

	bookeepers, err := generateAccounts(param.NodeIndexList)
	if err != nil {
		log.Errorf("failed to generate bookeepers, err: %v", err)
		return false
	}

	for _, bookeeper := range bookeepers {
		if _, err := bookeeper.SyncGenesisHeader(param.CrossChainID, rawHeader); err != nil {
			log.Errorf("failed to sync genesis header, err: %v", err)
			return false
		}
	}
	return true
}

func Mint() bool {
	var param struct {
		SideChainID  uint64
		CrossChainID uint64
		SideChainUrl string
		NodeKey      string
		Amount       uint64
	}

	log.Info("start to mint...")

	if err := config.LoadParams("test_mint.json", &param); err != nil {
		log.Errorf("failed to load params, err: %v", err)
		return false
	}

	amount := new(big.Int).Mul(ETH1, new(big.Int).SetUint64(param.Amount))
	sender, err := masterAccount()
	if err != nil {
		log.Errorf("failed to generate sender account, err: %v", err)
		return false
	}

	receiver, err := customGenerateAccount(param.SideChainUrl, param.SideChainID, param.NodeKey)
	if err != nil {
		log.Errorf("failed to generate receiver account, err: %v", err)
		return false
	}

	if _, err := sender.Mint(param.CrossChainID, receiver.Address(), amount); err != nil {
		log.Errorf("failed to mint token, err: %v", err)
	}

	// convert poly notify events and commitProof to side chain
	return true
}

func Burn() bool {
	return true
}
