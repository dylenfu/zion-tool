package core

import (
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/contracts/native/utils"
	"github.com/ethereum/go-ethereum/rlp"

	"github.com/dylenfu/zion-tool/config"
	"github.com/dylenfu/zion-tool/pkg/sdk"
	"github.com/ethereum/go-ethereum/common/hexutil"
	. "github.com/ethereum/go-ethereum/contracts/native/go_abi/node_manager_abi"
	nm "github.com/ethereum/go-ethereum/contracts/native/governance/node_manager"
	"github.com/ethereum/go-ethereum/crypto"
)

func masterAccount() (*sdk.Account, error) {
	chainID := config.Conf.ChainID
	node := config.Conf.Nodes[0]
	return sdk.CustomNewAccount(chainID, node.Url, node.PrivateKey)
}

func generateAccount(index int) (*sdk.Account, error) {
	chainID := config.Conf.ChainID
	if index >= len(config.Conf.Nodes) {
		return nil, fmt.Errorf("node index out of range")
	}
	node := config.Conf.Nodes[index]
	return sdk.CustomNewAccount(chainID, node.Url, node.PrivateKey)
}

func generateAccounts(indexList []int) ([]*sdk.Account, error) {
	chainID := config.Conf.ChainID
	nodes := config.Conf.Nodes
	list := make([]*sdk.Account, 0)
	for _, index := range indexList {
		if index >= len(nodes) {
			return nil, fmt.Errorf("node index out of range")
		}
		node := nodes[index]
		acc, err := sdk.CustomNewAccount(chainID, node.Url, node.PrivateKey)
		if err != nil {
			return nil, err
		}
		list = append(list, acc)
	}
	return list, nil
}

func getPeers(nodeIndexList []int) (*nm.Peers, error) {
	list := make([]*nm.PeerInfo, 0)
	nodes := config.Conf.Nodes
	for _, index := range nodeIndexList {
		if index >= len(nodes) {
			return nil, fmt.Errorf("node index out of range")
		}
		node := nodes[index]
		pubkey := hexutil.Encode(crypto.CompressPubkey(node.PublicKey))
		list = append(list, &nm.PeerInfo{PubKey: pubkey, Address: node.Address})
	}
	return &nm.Peers{List: list}, nil
}

func getProposalReceipt(proposer *sdk.Account, tx common.Hash) (*nm.EpochInfo, error) {
	receipt, err := proposer.GetReceipt(tx)
	if err != nil {
		return nil, err
	}
	event := receipt.Logs[0].Data
	list, err := utils.UnpackEvent(*nm.ABI, EventProposed, event)
	if err != nil {
		return nil, err
	}

	dec := list[0].([]byte)
	var epoch *nm.EpochInfo
	if err := rlp.DecodeBytes(dec, &epoch); err != nil {
		return nil, err
	}
	return epoch, nil
}
