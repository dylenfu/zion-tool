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
	"time"

	"github.com/dylenfu/zion-tool/config"
	"github.com/dylenfu/zion-tool/pkg/log"
	"github.com/ethereum/go-ethereum/common"
)

func Transfer() bool {
	var param struct {
		To     string
		Amount uint64
	}

	if err := config.LoadParams("test_transfer.json", &param); err != nil {
		log.Errorf("failed to load params, err: %v", err)
		return false
	}

	acc, err := masterAccount()
	if err != nil {
		log.Errorf("generate master account failed, err: %v", err)
	}
	to := common.HexToAddress(param.To)
	amount := new(big.Int).Mul(ETH1, new(big.Int).SetUint64(param.Amount))

	balanceBeforeTransfer, err := acc.BalanceOf(to, nil)
	if err != nil {
		log.Errorf("failed to get balance before transfer, err: %v", err)
		return false
	} else {
		log.Infof("balance before transfer %s", balanceBeforeTransfer.String())
	}

	if _, err := acc.Transfer(to, amount); err != nil {
		log.Errorf("failed to transfer eth, err: %v", err)
		return false
	}
	time.Sleep(5 * time.Second)

	balanceAfterTransfer, err := acc.BalanceOf(to, nil)
	if err != nil {
		log.Errorf("failed to get balance before transfer, err: %v", err)
		return false
	} else {
		log.Infof("balance after transfer %s", balanceAfterTransfer.String())
	}

	if balanceAfterTransfer.Cmp(new(big.Int).Add(balanceBeforeTransfer, amount)) != 0 {
		log.Error("balance not match")
		return false
	}

	return true
}
