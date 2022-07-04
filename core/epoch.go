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
	"math/big"

	"github.com/dylenfu/zion-tool/config"
	"github.com/dylenfu/zion-tool/pkg/log"
)

func Register() bool {
	var param struct {
		NodeIndexList []int
		StakeAmount   int
	}

	if err := config.LoadParams("test_register.json", &param); err != nil {
		log.Errorf("failed to load params, err: %v", err)
		return false
	}

	log.Split("start to change epoch")

	log.Split("start to prepare balance")
	if err := prepareBalance(); err != nil {
		log.Errorf("failed to prepare balance, err: %v", err)
		return false
	}

	vals, err := generateAccounts(param.NodeIndexList)
	if err != nil {
		log.Errorf("failed to generate proposer, err: %v", err)
		return false
	}

	log.Split("start to register nodes")
	stakeAmt := new(big.Int).Mul(big.NewInt(int64(param.StakeAmount)), ETH1)
	for _, v := range vals {
		balance, err := v.Balance(nil)
		if err != nil {
			return false
		} else {
			log.Infof("stake account %s balance %v", v.Addr().Hex(), balance)
		}
		if _, err := v.Register(v.PublicKey, stakeAmt, v.StakeAddr.Hex()); err != nil {
			log.Errorf("failed to register account, err: %v", err)
			return false
		}
	}

	wait()

	return true
}

func Stake() bool {
	return true
}
