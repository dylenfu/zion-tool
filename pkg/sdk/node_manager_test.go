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
	"crypto/ecdsa"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/stretchr/testify/assert"
	"math/big"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	. "github.com/ethereum/go-ethereum/contracts/native/go_abi/node_manager_abi"
	nm "github.com/ethereum/go-ethereum/contracts/native/governance/node_manager"
	"github.com/ethereum/go-ethereum/contracts/native/utils"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/rlp"
)

func nodeKeys() []string {
	return []string{
		"4b0c9b9d685db17ac9f295cb12f9d7d2369f5bf524b3ce52ce424031cafda1ae",
		"cc69b13ca2c5cd4d76bb881f6ad18d93bd947042c0f3a7adc80bdd17dac68210",
		"49e26aa4d60196153153388a24538c2693d65f0010a3a488c0c4c2b2a64b2de4",
		"9fc1723cff3bc4c11e903a53edb3b31c57b604bfc88a5d16cfec6a64fbf3141c",
		"5555ebb339d3d5ed1efbf0ca96f5b145134e5ce8044fec693558056d268776ae",
		"3d9c828244d3b2da70233a0a2aea7430feda17bded6edd7f0c474163802a431c",
		"018c71d5e3b245117ffba0975e46129371473c6a1d231c5eddf7a8364d704846",
		"c8d3e5e3fbc72898d1b90dedff34d6043fcbaaadeecd0bcb211a05c7c9a33af7",
	}
}

func getPeers(t *testing.T) []*nm.PeerInfo {
	keys := nodeKeys()
	list := make([]*nm.PeerInfo, 0)
	for _, v := range keys {
		pk, _ := crypto.HexToECDSA(v)
		pubkeyEnc := crypto.CompressPubkey(&pk.PublicKey)
		pubkey := hexutil.Encode(pubkeyEnc)
		addr := crypto.PubkeyToAddress(pk.PublicKey)
		list = append(list, &nm.PeerInfo{PubKey: pubkey, Address: addr})
	}
	return list
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
	t.Log(epoch.String())
}

// go test -v -count=1 github.com/dylenfu/zion-tool/pkg/sdk -run TestPropose
func TestPropose(t *testing.T) {
	testUrl = "http://localhost:22000"
	val1, _ := CustomNewAccount(testChainID, testUrl, getPrivateKey(0))

	startHeight := uint64(100)
	peers := &nm.Peers{List: getPeers(t)[0:5]}

	tx, err := val1.Propose(startHeight, peers)
	if err != nil {
		t.Fatal(err)
	} else {
		t.Logf("validator %s propose, hash %s", val1.Address().Hex(), tx.Hex())
	}

	t.Log("---------------------------------------------------------")
	t.Log("dump epoch")
	t.Log("---------------------------------------------------------")
	receipt, err := val1.GetReceipt(tx)
	if err != nil {
		t.Fatal(err)
	}
	event := receipt.Logs[0].Data
	list, err := utils.UnpackEvent(*nm.ABI, EventProposed, event)
	if err != nil {
		t.Fatal(err)
	}
	dec := list[0].([]byte)
	var epoch *nm.EpochInfo
	if err := rlp.DecodeBytes(dec, &epoch); err != nil {
		t.Fatal(err)
	}
	t.Logf(epoch.String())
}

// go test -v -count=1 github.com/dylenfu/zion-tool/pkg/sdk -run TestVote
func TestVote(t *testing.T) {
	testUrl = "http://localhost:22000"

	voters := getTestAccounts(3)
	hash := common.HexToHash("0xb8380e3573ca36c1d9b7f245da3e06e68038ef86d54708b5e348a1371f97e90f")

	curEpoch, err := voters[0].Epoch()
	if err != nil {
		t.Fatal(err)
	}
	epochID := curEpoch.ID + 1
	for _, voter := range voters {
		if tx, err := voter.Vote(epochID, hash); err != nil {
			t.Fatal(err)
		} else {
			t.Logf("voter %s voted, hash %s", voter.Address().Hex(), tx.Hex())
		}
		time.Sleep(5 * time.Second)
	}
}

func TestCommittedSeals(t *testing.T) {
	testUrl = "http://101.32.99.70:22000"

	blockNum := uint64(501)
	acc := getTestAccount(0)
	header, err := acc.client.HeaderByNumber(context.Background(), new(big.Int).SetUint64(blockNum))
	assert.NoError(t, err)

	extra, err := types.ExtractHotstuffExtraPayload(header.Extra)
	assert.NoError(t, err)
	t.Logf("extra committed seals size %d", len(extra.CommittedSeal))
}