package core

import (
	"fmt"
	"math/big"
	"time"

	"github.com/dylenfu/zion-tool/config"
	"github.com/dylenfu/zion-tool/pkg/log"
	"github.com/dylenfu/zion-tool/pkg/sdk"
	"github.com/ethereum/go-ethereum/params"
)

var ETH1 = params.ZNT1

type Account struct {
	*config.Node
	*sdk.Account
}

func masterAccount() (*sdk.Account, error) {
	chainID := config.Conf.ChainID
	node := config.Conf.Nodes[0]
	return sdk.CustomNewAccount(chainID, node.Url, node.PrivateKey)
}

func generateAccounts(indexList []int) ([]*Account, error) {
	list := make([]*Account, 0)
	for _, index := range indexList {
		if index >= len(config.Conf.Nodes) {
			return nil, fmt.Errorf("node index out of range")
		}
		acc, err := generateAccount(index)
		if err != nil {
			return nil, err
		}
		list = append(list, acc)
	}
	return list, nil
}

func generateAccount(index int) (*Account, error) {
	chainID := config.Conf.ChainID
	if index >= len(config.Conf.Nodes) {
		return nil, fmt.Errorf("node index out of range")
	}
	node := config.Conf.Nodes[index]
	acc, err := sdk.CustomNewAccount(chainID, node.Url, node.StakePrivateKey)
	if err != nil {
		return nil, err
	}
	return &Account{
		Node:    node,
		Account: acc,
	}, nil
}

func prepareBalance() error {
	amount := big.NewInt(int64(config.Conf.InitBalance))
	master, err := masterAccount()
	if err != nil {
		return err
	}

	// the first one is master account
	for i := 0; i < len(config.Conf.Nodes); i++ {
		addr := config.Conf.Nodes[i].StakeAddr
		balance, err := master.BalanceOf(addr, nil)
		if err != nil {
			return err
		}
		if balance.Cmp(amount) >= 0 {
			continue
		}
		added := new(big.Int).Sub(amount, balance)
		if _, err := master.Transfer(addr, added); err != nil {
			return err
		} else {
			log.Infof("prepare balance %v, added %v", amount, added)
		}
	}
	return nil
}

func wait() {
	time.Sleep(config.Conf.BlockWaitingTime())
}
