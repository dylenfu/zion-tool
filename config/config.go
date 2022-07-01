/*
 * Copyright (C) 2020 The poly network Authors
 * This file is part of The poly network library.
 *
 * The  poly network  is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Lesser General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * The  poly network  is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU Lesser General Public License for more details.
 * You should have received a copy of the GNU Lesser General Public License
 * along with The poly network .  If not, see <http://www.gnu.org/licenses/>.
 */

package config

import (
	"crypto/ecdsa"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strings"
	"time"

	"github.com/dylenfu/zion-tool/pkg/files"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
)

var (
	Conf *Config
)

type Config struct {
	Workspace   string
	ChainID     uint64
	Nodes       []*Node
	BlockPeriod int
	InitBalance int
}

func (c *Config) BlockWaitingTime() time.Duration {
	return time.Second * time.Duration(c.BlockPeriod+1)
}

type Node struct {
	NodeKey         string            `json:"NodeKey"`
	Url             string            `json:"Url"`
	StakeKey        string            `json:"StakeKey"`
	Address         common.Address    `json:"Address,omitempty"`
	PrivateKey      *ecdsa.PrivateKey `json:"PrivateKey,omitempty"`
	PublicKey       *ecdsa.PublicKey  `json:"PublicKey,omitempty"`
	StakeAddr       common.Address    `json:"StakeAddr,omitempty"`
	StakePrivateKey *ecdsa.PrivateKey `json:"StakePrivateKey,omitempty"`
	StakePublicKey  *ecdsa.PublicKey  `json:"StakePublicKey,omitempty"`
}

func LoadConfig(filepath string) {
	var (
		data []byte
		err  error
	)

	if data, err = files.ReadFile(filepath); err != nil {
		panic(err)
	}
	if err = json.Unmarshal(data, &Conf); err != nil {
		panic(err)
	}

	for index, v := range Conf.Nodes {
		v.PrivateKey, v.PublicKey, v.Address, err = ParsePrivateHex(v.NodeKey)
		if err != nil {
			panic(fmt.Sprintf("node key invalid, index %d, err: %v", index, err))
		}

		v.StakePrivateKey, v.StakePublicKey, v.StakeAddr, err = ParsePrivateHex(v.StakeKey)
		if err != nil {
			panic(fmt.Sprintf("stake key invalid, index %d, err: %v", index, err))
		}
	}
}

func LoadParams(fileName string, data interface{}) error {
	filepath := files.FullPath(Conf.Workspace, "cases", fileName)
	bz, err := ioutil.ReadFile(filepath)
	if err != nil {
		return err
	}
	return json.Unmarshal(bz, data)
}

func ParsePrivateHex(data string) (*ecdsa.PrivateKey, *ecdsa.PublicKey, common.Address, error) {
	key := data
	if !strings.Contains(key, "0x") {
		key = "0x" + key
	}

	enc, err := hexutil.Decode(key)
	if err != nil {
		return nil, nil, common.Address{}, err
	}

	pk, err := crypto.ToECDSA(enc)
	if err != nil {
		return nil, nil, common.Address{}, err
	}
	return pk, &pk.PublicKey, crypto.PubkeyToAddress(pk.PublicKey), nil
}
