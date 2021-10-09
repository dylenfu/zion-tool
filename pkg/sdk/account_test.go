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

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

var (
	testUrl                = "http://49.234.146.144:8545"
	testMainNodeKey        = "4b0c9b9d685db17ac9f295cb12f9d7d2369f5bf524b3ce52ce424031cafda1ae"
	testChainID     uint64 = 1002
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
	t.Logf("master %s transfer %v to %s, hash %s", master.Address(), amount, to.Hex(), hash.Hex())
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
