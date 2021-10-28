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

package config

import (
	"io/ioutil"
	"testing"

	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
)

func TestLoadConfig(t *testing.T) {
	filepath := "./dev.json"
	LoadConfig(filepath)
	for _, v := range Conf.Nodes {
		t.Log(v.Address.Hex())
	}
}

func TestLoadPrivateKey(t *testing.T) {
	files := []string{
		"/Users/dylen/software/hotstuff/zion-poa-dev/setup/keystore/node0/UTC--2021-10-28T06-11-29.578784000Z--e4113f21494aae440e551052a2b8666e9f33eaf2",
		"/Users/dylen/software/hotstuff/zion-poa-dev/setup/keystore/node1/UTC--2021-10-28T06-11-45.473969000Z--003deb81d33ea23254ec8c7727e0db6716c444cc",
		"/Users/dylen/software/hotstuff/zion-poa-dev/setup/keystore/node2/UTC--2021-10-28T06-11-52.377974000Z--dead7c3c6ae4869dc94158650e4cb3dd1eaaa718",
		"/Users/dylen/software/hotstuff/zion-poa-dev/setup/keystore/node3/UTC--2021-10-28T06-11-59.378911000Z--52f93cdc8a7085384ab0abe0433dbb2e0f16965a",
	}
	pwd := "111111"

	for _, file := range files {
		enc, err := ioutil.ReadFile(file)
		if err != nil {
			t.Fatal(err)
		}
		key, err := keystore.DecryptKey(enc, pwd)
		if err != nil {
			t.Fatal(err)
		}
		blob := crypto.FromECDSA(key.PrivateKey)
		hex := hexutil.Encode(blob)
		t.Logf("%s %s", key.Address.Hex(), hex)
	}
}
