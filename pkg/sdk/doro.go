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
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"

	"github.com/dylenfu/zion-tool/pkg/go_abi/doro"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/contracts/native/utils"
)

func (c *Account) SetDoro(contract common.Address, num uint64) (common.Hash, error) {
	ab, err := abi.JSON(strings.NewReader(doro.DoroABI))
	if err != nil {
		return common.EmptyHash, err
	}
	payload, err := utils.PackMethod(&ab, "setDoro", num)
	if err != nil {
		return common.EmptyHash, err
	}
	return c.signAndSendTx(payload, contract)
}

func (c *Account) makeAuthWithoutGasLimit() (*bind.TransactOpts, error) {
	fromAddress := c.Addr()
	nonce, err := c.client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		return nil, err
	}

	gasPrice, err := c.client.SuggestGasPrice(context.Background())
	if err != nil {
		return nil, err
	}

	auth := bind.NewKeyedTransactor(c.pk)
	auth.Nonce = big.NewInt(int64(nonce))
	auth.Value = big.NewInt(int64(0)) // in wei
	//auth.GasLimit = uint64(0) // in units
	auth.GasPrice = gasPrice.Mul(gasPrice, big.NewInt(1))
	return auth, nil
}
