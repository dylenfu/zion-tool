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
	"fmt"
	"io/ioutil"

	"github.com/dylenfu/zion-tool/config"
	"github.com/dylenfu/zion-tool/pkg/log"
	"github.com/dylenfu/zion-tool/pkg/sdk"
	"github.com/ethereum/go-ethereum/accounts/keystore"
)

func NeoProof() bool {
	var params struct {
		Keystore string
		Key      string
		Value    string
		LastTime int
	}

	if err := config.LoadParams("test_proof.json", &params); err != nil {
		log.Errorf("failed to load params json file, err: %v", err)
		return false
	}

	chainID := config.Conf.ChainID
	url := config.Conf.Nodes[0].Url
	enc, err := ioutil.ReadFile(params.Keystore)
	if err != nil {
		log.Errorf("failed to load keystore json file, err: %v", err)
		return false
	}

	pk, err := keystore.DecryptKey(enc, "111111")
	if err != nil {
		log.Errorf("failed to decrypt ecdsa key, err: %v", err)
		return false
	}

	acc, err := sdk.CustomNewAccount(chainID, url, pk.PrivateKey)
	if err != nil {
		log.Errorf("failed to generate account, err: %v", err)
		return false
	}

	for i := 0; i < params.LastTime; i++ {
		if tx, err := acc.NeoProof(params.Key, params.Value); err != nil {
			log.Errorf("failed to set proof, err: %v", err)
			return false
		} else {
			log.Split(fmt.Sprintf("index %d, tx %s: set proof (%s, %s)", i, tx.Hex(), params.Key, params.Value))
		}
	}

	return true
}
