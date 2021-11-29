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
	"math/big"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	. "github.com/ethereum/go-ethereum/contracts/native/go_abi/node_manager_abi"
	nm "github.com/ethereum/go-ethereum/contracts/native/governance/node_manager"
	"github.com/ethereum/go-ethereum/contracts/native/utils"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/stretchr/testify/assert"
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

/*
neo machine
alias ssh-zion-neo1='ssh ubuntu@49.234.146.144 -p 32000'
alias ssh-zion-neo2='ssh ubuntu@42.192.185.27 -p 32000'
alias ssh-zion-neo3='ssh ubuntu@212.129.153.164 -p 32000'
alias ssh-zion-neo4='ssh ubuntu@42.192.186.75 -p 32000'
*/
func TestGetProof(t *testing.T) {
	urls := map[string]string{
		"node0": "http://49.234.146.144:8545",
		"node1": "http://42.192.185.27:8545",
		"node2": "http://212.129.153.164:8545",
		"node3": "http://42.192.186.75:8545",
		"node4": "http://49.234.146.144:8645",
		"node5": "http://42.192.185.27:8645",
		"node6": "http://212.129.153.164:8645",
	}
	cid := uint64(1002)
	blockNum := uint64(93442)

	contract := common.HexToAddress("0xa6063efed0487cf5d519b0c9e298c7d2da82e74e")
	//keys := []common.Hash{
	//	common.HexToHash("0xa8e2f95eb5af99bc72e1b4d25c76d5abc83861cb2401fcd920d4a5249b1246a6"),
	//	common.HexToHash("0x976b3f7bcb9b91ec9493a20aafb2fcf89745ca105ee2b1921c72037ceac1a1a8"),
	//}
	//for node, url := range urls {
	//	acc, err := NewAccount(cid, url)
	//	assert.NoError(t, err)
	//	for _, key := range keys {
	//		enc, err := acc.client.StorageAt(context.Background(), contract, key, new(big.Int).SetUint64(blockNum))
	//		assert.NoError(t, err)
	//		t.Logf("%s storage state (%s, %s) at block %d", node, key.Hex(), hexutil.Encode(enc), blockNum)
	//	}
	//}
	//
	//acc, err := NewAccount(cid, "http://212.129.153.164:8545")
	//assert.NoError(t, err)
	//
	//for _, key := range keys {
	//	blockNum += 1
	//	enc, err := acc.client.StorageAt(context.Background(), contract, key, new(big.Int).SetUint64(blockNum))
	//	if err == nil {
	//		t.Logf("node3 storage state (%s, %s) at block %d", key.Hex(), hexutil.Encode(enc), blockNum)
	//	}
	//}

	pkeys := []string{
		"0xa8e2f95eb5af99bc72e1b4d25c76d5abc83861cb2401fcd920d4a5249b1246a6",
		"0x976b3f7bcb9b91ec9493a20aafb2fcf89745ca105ee2b1921c72037ceac1a1a8",
	}
	for node, url := range urls {
		acc, _ := NewAccount(cid, url)
		res, err := acc.client.ProofAt(context.Background(), contract, pkeys, new(big.Int).SetUint64(blockNum))
		assert.NoError(t, err)
		t.Logf("%s proof %v", node, res)
	}
}

func TestEstimate(t *testing.T) {
	urls := map[string]string{
		//"node0": "http://49.234.146.144:8545",
		//"node1": "http://42.192.185.27:8545",
		"node2": "http://212.129.153.164:8545",
		//"node3": "http://42.192.186.75:8545",
		//"node4": "http://49.234.146.144:8645",
		//"node5": "http://42.192.185.27:8645",
		//"node6": "http://212.129.153.164:8645",
	}
	cid := uint64(1002)
	blockNum := uint64(93442)
	//
	//pkeys := []string{
	//	"0xa8e2f95eb5af99bc72e1b4d25c76d5abc83861cb2401fcd920d4a5249b1246a6",
	//	"0x976b3f7bcb9b91ec9493a20aafb2fcf89745ca105ee2b1921c72037ceac1a1a8",
	//}

	from := common.HexToAddress("0x67cde763bd045b14898d8b044f8afc8695ae8608")
	txhash := common.HexToHash("0x5d6db1fa75b36d55773f741a090bcb5da1f1b0f0912e443dcc88ec47fd963ff6")

	acc2, _ := NewAccount(cid, urls["node2"])
	tx, _, err := acc2.client.TransactionByHash(context.Background(), txhash)
	assert.NoError(t, err)

	for _, url := range urls {
		acc, _ := NewAccount(cid, url)
		acc.CustomEstimate(tx, from, new(big.Int).SetUint64(blockNum))
		//acc.client.EstimateGas(context.Background(), ethereum.CallMsg{})
		//estimate(acc, txhash, from)
		//res, err := acc.client.ProofAt(context.Background(), contract, pkeys, new(big.Int).SetUint64(blockNum))
		//assert.NoError(t, err)
		//t.Logf("%s proof %v", node, res)
	}
}

// go test -count=1 -v github.com/dylenfu/zion-tool/pkg/sdk -run TestEstimateSyncBlockHeader
func TestEstimateSyncBlockHeader(t *testing.T) {
	chainID := uint64(102)
	crossChainID := uint64(10002)
	hexBlob := "0x7b22706172656e7448617368223a22307865326637393932666637353035373830353333343238393564383662643361383134353738393564616135633137343935366434323261636262663763303538222c2273686133556e636c6573223a22307831646363346465386465633735643761616238356235363762366363643431616433313234353162393438613734313366306131343266643430643439333437222c226d696e6572223a22307865396537303334616564356365376635623064323831636665333437623861356332633533353034222c227374617465526f6f74223a22307834303733636335373938386631313337663732393162643334656233346534356630336166316431333533323136306661333735653832616435303533623162222c227472616e73616374696f6e73526f6f74223a22307839633336663834613361646431343430343861633162306466303537663035306335643964303236656261656530333831396530363932616436386535353864222c227265636569707473526f6f74223a22307866626264306437393530636365393439393639383039646364343964633166323562653437393131663333336438383963376537393037303635333330653033222c226c6f6773426c6f6f6d223a2230783030323030303230303130303030303030303030303030303830303030303030303030303030303034303030303031303030303030313030303030303530303030303430303030303034303030343030303030303030303030303030303030303030303030303030303030303030303138303030303030303030303030303030303430303030303030303030303030303230303030303038383030303030323030383030303030383030303030303030303030303030303038303830303030303030303030303030313030303030303030303030303038303030303030303030303030303030303130303030303030303030303030303130303030323030303030323030303230303030303030303030303030303030303030303030303230323030303030303031303038303230303830303030303034303030303030303030303030303030303030303030303031303230303030303030303230303038303030303430303030303030303030303130303030303030303030303030303830303132303030303032303034313830303030303030303030303030303030303130303030303030303130303030303039303030303030303030303030303030303030303030303030303030303030303030303030303030303030383032303030313030303032303230303030303030343030303030303030313830303030303030222c22646966666963756c7479223a2230783839333339323631222c226e756d626572223a223078616435333937222c226761734c696d6974223a223078376131323030222c2267617355736564223a223078313361363131222c2274696d657374616d70223a2230783631383338326666222c22657874726144617461223a22307836333733326537303634373832653635363437353230363736353734363832303736333132653331333032653338222c226d697848617368223a22307863613965666463376338333364343035636435356365313635353937623866323739366361346636356438313766643962656433303030333234326466303466222c226e6f6e6365223a22307836373063653739326138633663353963222c2262617365466565506572476173223a22307862222c2268617368223a22307838383232643638356263333237383464363030653230333231633463303338373939626161323564623931393465613333373739303037346131353462366131227d"
	url := "http://127.0.0.1:22007"

	raw, err := hexutil.Decode(hexBlob)
	if err != nil {
		t.Fatal(err)
	}
	headers := [][]byte{raw}

	acc, _ := CustomNewAccount(chainID, url, getPrivateKey(0))

	gas, err := acc.EstimateSubmitHeaders(crossChainID, headers)
	if err != nil {
		t.Fatal(err)
	} else {
		t.Logf("estimate gas %d", gas)
	}
}
