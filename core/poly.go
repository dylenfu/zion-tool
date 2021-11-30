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
	"encoding/hex"
	scom "github.com/ethereum/go-ethereum/contracts/native/cross_chain_manager/common"
	"github.com/ethereum/go-ethereum/core/state"
	"math/big"

	"github.com/dylenfu/zion-tool/config"
	"github.com/dylenfu/zion-tool/pkg/log"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/contracts/native/governance/node_manager"
	"github.com/ethereum/go-ethereum/contracts/native/utils"
	"github.com/ethereum/go-ethereum/rlp"
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

// sync side chain genesis header to poly
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

// 查询poly上最近一次epoch内容， eccm中初始化genesisBlock或者changeEpoch时，假如：
// 1. block99的header记录下一组epoch(ep2)valset
// 2. block100正式切换epoch，新的valset生效，并参与共识
// 3. 提交给eccm的参数是header100, ep2, 而不是header99，ep2.
//    因为eccm中是根据epoch携带的valset做变更，而不是header中携带的valset做变更
// 4.
func FetchEpochProof() bool {
	var param struct {
		Height  uint64
		EpochId uint64
		IsGenesis bool
	}

	log.Info("start to fetch epoch header and proof...")

	if err := config.LoadParams("test_epoch_proof.json", &param); err != nil {
		log.Errorf("failed to load params, err: %v", err)
		return false
	}

	// param height should be epoch start height, and block n- 1 is the block which stored new epoch validators
	blockHeight := param.Height
	if !param.IsGenesis {
		blockHeight = param.Height- 1
	}

	sdk, err := masterAccount()
	if err != nil {
		log.Errorf("failed to load master account, err: %v", err)
		return false
	}

	// get header and raw committed seals
	_, rawHeader, rawSeals, err := sdk.GetRawHeaderAndSeals(blockHeight)
	log.Infof("rawHeader: %s", hexutil.Encode(rawHeader))
	log.Infof("rawSeal: %s", hexutil.Encode(rawSeals))

	// get epoch
	epoch, err := sdk.GetEpochByID(param.EpochId, "latest")
	if err != nil {
		log.Errorf("failed to load epoch, err: %v", err)
		return false
	}
	var inf = struct {
		ID          uint64
		Peers       *node_manager.Peers
		StartHeight uint64
	}{
		ID:          epoch.ID,
		Peers:       epoch.Peers,
		StartHeight: epoch.StartHeight,
	}

	// get proof
	contractAddr := utils.NodeManagerContractAddress
	proofHash := node_manager.EpochProofHash(epoch.ID)
	cacheKey := utils.ConcatKey(contractAddr, []byte(node_manager.SKP_PROOF), proofHash.Bytes())
	slot := state.Key2Slot(cacheKey[common.AddressLength:])
	key := hexutil.Encode(slot[:])
	storageKeys := []string{key}
	blockNum := new(big.Int).SetUint64(epoch.StartHeight)
	accountProof, storageProof, err := sdk.GetProof(contractAddr, storageKeys, blockNum)
	if err != nil {
		log.Errorf("failed to get proof, err: %v", err)
	}
	log.Infof("account proof: %s", hexutil.Encode(accountProof))
	log.Infof("storage proof: %s", hexutil.Encode(storageProof))

	// raw epoch
	rawEpoch, err := rlp.EncodeToBytes(inf)
	if err != nil {
		log.Errorf("failed to encode epoch, err: %v", err)
		return false
	}
	log.Infof("rawEpoch: %s", hexutil.Encode(rawEpoch))

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

	// mint token on main chain
	hash, err := sender.Mint(param.CrossChainID, receiver.Address(), amount)
	if err != nil {
		log.Errorf("failed to mint token, err: %v", err)
		return false
	} else {
		log.Infof("mint success, hash %s", hash.Hex())
	}

	// fetch receipt on main chain
	receipt, err := sender.GetReceipt(hash)
	if err != nil {
		log.Errorf("failed to get receipt, err: %v", err)
		return false
	}

	// assemble params for `verifyHeaderAndExecuteTx`
	_, rawHeader, rawSeals, err := sender.GetRawHeaderAndSeals(receipt.BlockNumber.Uint64())
	if err != nil {
		log.Errorf("failed to get raw header and seals, err: %v", err)
		return false
	}

	log.Infof("rawHeader: %s", hexutil.Encode(rawHeader))
	log.Infof("rawSeals: %s", hexutil.Encode(rawSeals))

	// get proof
	if len(receipt.Logs) != 3 {
		log.Errorf("event logs should contain (crossChainEvent, lockEvent, makeProofNotify)")
	}
	notify := receipt.Logs[2]
	list, err := utils.UnpackEvent(*scom.ABI, scom.NOTIFY_MAKE_PROOF_EVENT, notify.Data)
	if err != nil {
		log.Errorf("failed to unpack makeProof, err: %v", err)
		return false
	}
	if len(list) != 3 {
		log.Errorf("unpacked list length != 3, it should contains crossChain, lock and makeProof tx event")
		return false
	}
	rawMerkelValue := list[0].(string)
	merkelDec, err := hex.DecodeString(rawMerkelValue)
	if err != nil {
		log.Errorf("failed to decode raw merkel value, err: %v", err)
		return false
	}
	log.Infof("merkle value: %s", hexutil.Encode(merkelDec))

	rawKey := list[2].(string)
	raw, err := hex.DecodeString(rawKey)
	if err != nil {
		log.Errorf("decode ")
	}
	slot := state.Key2Slot(raw[common.AddressLength:])
	key := hexutil.Encode(slot[:])
	storageKeys := []string{key}
	accountProof, storageProof, err := sender.GetProof(utils.CrossChainManagerContractAddress, storageKeys, receipt.BlockNumber)
	if err != nil {
		log.Errorf("failed to get proof, err: %v", err)
	}
	log.Infof("account proof: %s", hexutil.Encode(accountProof))
	log.Infof("storage proof: %s", hexutil.Encode(storageProof))

	return true
}

func Burn() bool {
	return true
}
