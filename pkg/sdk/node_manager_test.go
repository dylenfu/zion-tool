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
	"crypto/ecdsa"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common/hexutil"
	nm "github.com/ethereum/go-ethereum/contracts/native/governance/node_manager"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/params"
)

func nodeKeys() []string {
	return []string{
		"4b0c9b9d685db17ac9f295cb12f9d7d2369f5bf524b3ce52ce424031cafda1ae",
		"3d9c828244d3b2da70233a0a2aea7430feda17bded6edd7f0c474163802a431c",
		"cc69b13ca2c5cd4d76bb881f6ad18d93bd947042c0f3a7adc80bdd17dac68210",
		"018c71d5e3b245117ffba0975e46129371473c6a1d231c5eddf7a8364d704846",
		"49e26aa4d60196153153388a24538c2693d65f0010a3a488c0c4c2b2a64b2de4",
		"9fc1723cff3bc4c11e903a53edb3b31c57b604bfc88a5d16cfec6a64fbf3141c",
		"5555ebb339d3d5ed1efbf0ca96f5b145134e5ce8044fec693558056d268776ae",
		"a1a470badc2b949188796d5595be39b76249d840de231233667246ba1aca7613",
	}
}

func getPrivateKey(n int) *ecdsa.PrivateKey {
	keys := nodeKeys()
	pk, _ := crypto.HexToECDSA(keys[n])
	return pk
}

func getTestAccount(index int) *Account {
	acc, _ := CustomNewAccount(testChainID, testUrl, getPrivateKey(index))
	return acc
}

func getTestAccounts(n int) []*Account {
	list := make([]*Account, 0)
	for i := 0; i < n; i++ {
		list = append(list, getTestAccount(i))
	}
	return list
}

func TestGetEpoch(t *testing.T) {
	epoch, err := master.Epoch()
	if err != nil {
		t.Fatal(err)
	}
	t.Log("id", epoch.ID)
	t.Log("start", epoch.StartHeight)
	for index, v := range epoch.Validators {
		t.Log("index", index, v.Address.Hex())
	}
}

func TestRegister(t *testing.T) {
	amount := nm.GenesisMinInitialStake
	depositAmount := new(big.Int).Add(amount, params.ZNT1)
	stakePK, _ := crypto.GenerateKey()
	stakeAddr := crypto.PubkeyToAddress(stakePK.PublicKey)
	if _, err := master.Transfer(stakeAddr, depositAmount); err != nil {
		t.Error(err)
	}
	stakeAcc, _ := CustomNewAccount(testChainID, testUrl, stakePK)
	balance, err := stakeAcc.Balance(nil)
	if err != nil {
		t.Error(err)
	} else {
		t.Log("addr", stakeAcc.Addr().Hex(), "balance", balance)
	}
	if _, err := stakeAcc.Register(&master.pk.PublicKey, amount, "test1"); err != nil {
		t.Error(err)
	}
}

func TestStake(t *testing.T) {
	amount := nm.GenesisMinInitialStake
	stakePK, _ := crypto.GenerateKey()
	stakeAddr := crypto.PubkeyToAddress(stakePK.PublicKey)
	if _, err := master.Transfer(stakeAddr, amount); err != nil {
		t.Error(err)
	}
	stakeAcc, _ := CustomNewAccount(testChainID, testUrl, stakePK)
	if _, err := stakeAcc.Stake(&master.pk.PublicKey, amount); err != nil {
		t.Error(err)
	}
}

func TestGenerateKeys(t *testing.T) {
	for i := 0; i < 10; i++ {
		key, _ := crypto.GenerateKey()
		enc := crypto.FromECDSA(key)
		raw := hexutil.Encode(enc)
		t.Log(i, raw)
	}

	enc, _ := hexutil.Decode("0xa1a470badc2b949188796d5595be39b76249d840de231233667246ba1aca7613")
	pk, _ := crypto.ToECDSA(enc)
	raw := hexutil.Encode(crypto.CompressPubkey(&pk.PublicKey))
	addr := crypto.PubkeyToAddress(pk.PublicKey)
	t.Log(raw, addr)
}
