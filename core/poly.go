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
	"github.com/ethereum/go-ethereum/core/state"
	"math/big"

	"github.com/dylenfu/zion-tool/config"
	"github.com/dylenfu/zion-tool/pkg/log"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/contracts/native/governance/node_manager"
	"github.com/ethereum/go-ethereum/contracts/native/utils"
	"github.com/ethereum/go-ethereum/core/types"
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
	header, err := sdk.BlockHeaderByNumber(blockHeight)
	if err != nil {
		log.Errorf("failed to fetch header, err: %v", err)
		return false
	}
	rawHeader, err := rlp.EncodeToBytes(types.HotstuffFilteredHeader(header, false))
	if err != nil {
		log.Errorf("failed to encode filtered header")
		return false
	}
	extra, err := types.ExtractHotstuffExtra(header)
	if err != nil {
		log.Errorf("failed to extract header")
		return false
	}
	rawSeals, err := rlp.EncodeToBytes(extra.CommittedSeal)
	if err != nil {
		log.Errorf("failed to encode committed seals")
		return false
	}
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

	rlpEncodeStringList := func(raw []string) ([]byte, error) {
		var rawBytes []byte
		for i := 0; i < len(raw); i++ {
			rawBytes = append(rawBytes, common.Hex2Bytes(raw[i][2:])...)
		}
		return rlp.EncodeToBytes(rawBytes)
	}

	if blob, err := rlpEncodeStringList(accountProof); err != nil {
		log.Errorf("rlp encode account proof err: %v", err)
		return false
	} else {
		log.Infof("account proof: %s", hexutil.Encode(blob))
	}
	if blob, err := rlpEncodeStringList(storageProof); err != nil {
		log.Errorf("rlp encode storage proof err: %v", err)
		return false
	} else {
		log.Infof("storage proof: %s", hexutil.Encode(blob))
	}

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

	if hash, err := sender.Mint(param.CrossChainID, receiver.Address(), amount); err != nil {
		log.Errorf("failed to mint token, err: %v", err)
		return false
	} else {
		log.Infof("mint success, hash %s", hash.Hex())
	}

	// convert poly notify events and commitProof to side chain
	return true
}

func Burn() bool {
	return true
}
