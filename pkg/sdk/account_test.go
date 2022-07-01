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
	"os"
	"testing"

	"github.com/dylenfu/zion-tool/pkg/math"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/assert"
)

var (
	testUrl                = "http://127.0.0.1:22000"
	testMainNodeKey        = "4b0c9b9d685db17ac9f295cb12f9d7d2369f5bf524b3ce52ce424031cafda1ae"
	testChainID     uint64 = 60801
	master          *Account
	testEth1        = math.Pow10toBigInt(18)
)

func TestMain(m *testing.M) {
	privKey, _ := crypto.HexToECDSA(testMainNodeKey)
	master, _ = CustomNewAccount(testChainID, testUrl, privKey)
	os.Exit(m.Run())
}

// go test -v github.com/dylenfu/zion-tool/pkg/sdk -run TestTransfer
func TestTransfer(t *testing.T) {
	to := common.HexToAddress("0x67CDE763bD045B14898d8B044F8afC8695ae8608")
	amount := 1000000000
	value := new(big.Int).Mul(testEth1, new(big.Int).SetUint64(uint64(amount)))
	hash, err := master.Transfer(to, value)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("master %s transfer %v to %s, hash %s", master.Addr(), amount, to.Hex(), hash.Hex())
}

func TestGetBlock(t *testing.T) {
	start := 30
	end := 60
	for i := start; i < end; i++ {
		blockNum := big.NewInt(int64(i))
		block, err := master.client.BlockByNumber(context.Background(), blockNum)
		if err != nil {
			t.Fatal(err)
		}
		for _, tx := range block.Transactions() {
			t.Logf("tx %s, send to %s value %v", tx.Hash(), tx.To(), tx.Value())
		}
	}
}

func TestBlockHeaderRoot(t *testing.T) {
	height := uint64(222917)
	chainID := uint64(1002)
	nodesUrl := map[string]string{
		"node1": "http://49.234.146.144:8545",
		"node4": "http://49.234.146.144:8645",
		"node2": "http://42.192.185.27:8545",
		"node6": "http://42.192.185.27:8645",
		"node3": "http://212.129.153.164:8545",
		"node7": "http://212.129.153.164:8645",
		"node5": "http://42.192.186.75:8545",
	}
	accounts := make(map[string]*Account)
	for id, v := range nodesUrl {
		acc, _ := NewAccount(chainID, v)
		accounts[id] = acc
	}
	for id, v := range accounts {
		header, err := v.BlockHeaderByNumber(height)
		assert.NoError(t, err)
		t.Logf("%s header root %s", id, header.Root.Hex())
	}
}

func TestEstimateTx(t *testing.T) {
	chainID := uint64(1002)
	nodesUrl := map[string]string{
		"node1": "http://49.234.146.144:8545",
		"node4": "http://49.234.146.144:8645",
		"node2": "http://42.192.185.27:8545",
		"node6": "http://42.192.185.27:8645",
		"node3": "http://212.129.153.164:8545",
		"node7": "http://212.129.153.164:8645",
		"node5": "http://42.192.186.75:8545",
	}
	from := common.HexToAddress("0x67cde763bd045b14898d8b044f8afc8695ae8608")
	txhash := common.HexToHash("0x03d2cddfab5e19e2bd00cb81534ba67cb0a5674c1cab6d997871cc1fe85715ff")

	//chainID := uint64(102)
	//nodesUrl := map[string]string{
	//	"node1": "http://127.0.0.1:22000",
	//	"node2": "http://127.0.0.1:22001",
	//	"node3": "http://127.0.0.1:22002",
	//	"node4": "http://127.0.0.1:22003",
	//	"node5": "http://127.0.0.1:22004",
	//	"node6": "http://127.0.0.1:22005",
	//	"node7": "http://127.0.0.1:22006",
	//}
	//from := common.HexToAddress("0x258af48e28E4A6846E931dDfF8e1Cdf8579821e5")
	//txhash := common.HexToHash("0x481a9074cd7cc8182f2cd98d261a48f332cf9d3944369b3af4c3031d56f20e35")

	accounts := make(map[string]*Account)
	for id, v := range nodesUrl {
		acc, _ := NewAccount(chainID, v)
		accounts[id] = acc
	}
	for id, v := range accounts {
		if err := estimate(v, txhash, from); err != nil {
			t.Logf("%s: %s", id, err.Error())
		}
	}
}

func estimate(s *Account, txhash common.Hash, from common.Address) error {
	tx, _, err := s.client.TransactionByHash(context.Background(), txhash)
	if err != nil {
		return err
	}

	args := ethereum.CallMsg{
		From:       from,
		To:         tx.To(),
		Gas:        tx.Gas(),
		GasPrice:   tx.GasPrice(),
		Value:      tx.Value(),
		Data:       tx.Data(),
		AccessList: tx.AccessList(),
	}
	if _, err := s.client.EstimateGas(context.Background(), args); err != nil {
		return err
	}
	return nil
}

func (c *Account) CustomEstimate(tx *types.Transaction, from common.Address, blockNum *big.Int) error {
	args := ethereum.CallMsg{
		From:       from,
		To:         tx.To(),
		Gas:        tx.Gas(),
		GasPrice:   tx.GasPrice(),
		Value:      tx.Value(),
		Data:       tx.Data(),
		AccessList: tx.AccessList(),
	}
	if _, err := c.client.EstimateGas(context.Background(), args); err != nil {
		return err
	}
	return nil
}
