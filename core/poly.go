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
	"fmt"
	"math/big"
	"time"

	"github.com/dylenfu/zion-tool/pkg/sdk"
	scom "github.com/ethereum/go-ethereum/contracts/native/cross_chain_manager/common"
	"github.com/ethereum/go-ethereum/core/state"
	"github.com/ethereum/go-ethereum/core/types"

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
		Height    uint64
		EpochId   uint64
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
		blockHeight = param.Height - 1
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
		SideChainID   uint64
		CrossChainID  uint64
		SideChainUrl  string
		SideChainECCM string
		NodeKey       string
		Amount        uint64
		Relayer       bool
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

	senderBalanceBeforeMint, receiverBalanceBeforeMint, err := dumpBalance(sender, receiver, sender.Address(), receiver.Address(), true)
	if err != nil {
		log.Errorf("failed to get balance, err: %v", err)
		return false
	}

	// mint token on main chain
	log.Splitf("start to mint...")
	hash, err := sender.Mint(param.CrossChainID, receiver.Address(), amount)
	if err != nil {
		log.Errorf("failed to mint token, err: %v", err)
		return false
	}
	log.Splitf("mint success, hash %s", hash.Hex())

	// fetch receipt on main chain
	log.Splitf("start to fetch params...")
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
	accountProof, storageProof, merkelValue, err := transaction2Proof(sender, receipt, utils.CrossChainManagerContractAddress)
	if err != nil {
		log.Errorf("failed to get proof, err: %v", err)
		return false
	}
	log.Splitf("account proof: %s", hexutil.Encode(accountProof))
	log.Infof("storage proof: %s", hexutil.Encode(storageProof))
	log.Infof("merkel value: %s", hexutil.Encode(merkelValue))

	// return if we use relayer to commit proof
	if param.Relayer {
		return true
	}

	time.Sleep(3 * time.Second)
	hash, err = receiver.SideChainVerifyHeaderAndExecute(
		common.HexToAddress(param.SideChainECCM),
		rawHeader, rawSeals,
		accountProof, storageProof,
		merkelValue,
	)
	if err != nil {
		log.Errorf("failed to execute verifyHeaderAndExecute tx, err: %v", err)
		return false
	}
	log.Splitf("verifyHeaderAndExecuteTx success, tx hash %s", hash.Hex())

	time.Sleep(3 * time.Second)
	senderBalanceAfterMint, receiverBalanceAfterMint, err := dumpBalance(sender, receiver, sender.Address(), receiver.Address(), true)
	if err != nil {
		log.Errorf("failed to get balance, err: %v", err)
		return false
	}

	senderSubed := new(big.Int).Sub(senderBalanceBeforeMint, senderBalanceAfterMint)
	receiverAdded := new(big.Int).Sub(receiverBalanceAfterMint, receiverBalanceBeforeMint)
	if senderSubed.Cmp(amount) != 0 {
		log.Errorf("sender subed != amount, (%s, %s)", senderSubed.String(), amount.String())
		return false
	}
	if receiverAdded.Cmp(amount) != 0 {
		log.Errorf("receiver added != amount, (%s, %s)", receiverAdded.String(), amount.String())
		return false
	}
	return true
}

func Burn() bool {
	return true
}

func transaction2Proof(sender *sdk.Account, receipt *types.Receipt, contract common.Address) ([]byte, []byte, []byte, error) {
	if len(receipt.Logs) != 3 {
		log.Errorf("event logs should contain (crossChainEvent, lockEvent, makeProofNotify)")
	}
	notify := receipt.Logs[2]
	list, err := utils.UnpackEvent(*scom.ABI, scom.NOTIFY_MAKE_PROOF_EVENT, notify.Data)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to unpack makeProof, err: %v", err)
	}
	if len(list) != 3 {
		return nil, nil, nil, fmt.Errorf("unpacked list length != 3, it should contains crossChain, lock and makeProof tx event")
	}
	rawMerkelValue := list[0].(string)
	merkelDec, err := hex.DecodeString(rawMerkelValue)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to decode raw merkel value, err: %v", err)
	}

	rawKey := list[2].(string)
	raw, err := hex.DecodeString(rawKey)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("decode rawKey err: %v", err)
	}
	slot := state.Key2Slot(raw[common.AddressLength:])
	key := hexutil.Encode(slot[:])
	storageKeys := []string{key}
	accountProof, storageProof, err := sender.GetProof(utils.CrossChainManagerContractAddress, storageKeys, receipt.BlockNumber)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to get proof, err: %v", err)
	}
	return accountProof, storageProof, merkelDec, nil
}

func dumpBalance(mainChainSdk *sdk.Account, sideChainSdk *sdk.Account, sender, receiver common.Address, fromMainChain bool) (*big.Int, *big.Int, error) {
	senderBalanceOnMainChain, err := mainChainSdk.BalanceOf(sender, nil)
	if err != nil {
		return nil, nil, err
	}
	receiverBalanceOnMainChain, err := mainChainSdk.BalanceOf(receiver, nil)
	if err != nil {
		return nil, nil, err
	}

	senderBalanceOnSideChain, err := sideChainSdk.BalanceOf(sender, nil)
	if err != nil {
		return nil, nil, err
	}
	receiverBalanceOnSideChain, err := sideChainSdk.BalanceOf(receiver, nil)
	if err != nil {
		return nil, nil, err
	}

	log.Splitf("main chain balance:(sender, receiver) (%s, %s)", senderBalanceOnMainChain.String(), receiverBalanceOnMainChain.String())
	log.Infof("side chain balance:(sender, receiver) (%s, %s)", senderBalanceOnSideChain.String(), receiverBalanceOnSideChain.String())

	if fromMainChain {
		return senderBalanceOnMainChain, receiverBalanceOnSideChain, nil
	} else {
		return senderBalanceOnSideChain, receiverBalanceOnMainChain, nil
	}
}