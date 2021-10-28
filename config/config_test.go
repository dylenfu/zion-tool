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
		"/Users/dylen/software/hotstuff/zion-poa/setup/keystore/node0/UTC--2021-10-27T09-58-23.849917000Z--5ac410790b489c400594ad3a284141b4d0b38db5",
		"/Users/dylen/software/hotstuff/zion-poa/setup/keystore/node1/UTC--2021-10-27T09-58-30.071110000Z--b11772fb50cbfc63b4d853dd38c412867e4bf2f3",
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
