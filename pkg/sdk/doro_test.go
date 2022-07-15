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
	"testing"

	"github.com/dylenfu/zion-tool/pkg/go_abi/doro"
	"github.com/ethereum/go-ethereum/common"
)

func TestDoro1(t *testing.T) {
	acc := getStakeAccount()
	contract := common.HexToAddress("0x73b0727DA810d0be51D74E83655398fA6DC828aa")
	num := uint64(1234)
	if _, err := acc.SetDoro(contract, num); err != nil {
		t.Error(err)
	}
}

func TestDoro2(t *testing.T) {
	acc := getStakeAccount()
	contract := common.HexToAddress("0x73b0727DA810d0be51D74E83655398fA6DC828aa")
	num := uint64(12)

	instance, err := doro.NewDoro(contract, acc.client)
	if err != nil {
		t.Fatal(err)
	}
	auth, err := acc.makeAuthWithoutGasLimit()
	if err != nil {
		t.Fatal(err)
	}

	tx, err := instance.SetDoro(auth, num)
	if err != nil {
		t.Fatal(err)
	}
	if err := acc.WaitTransaction(tx.Hash()); err != nil {
		t.Fatal(err)
	}
	got, err := instance.Data(nil)
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("expect num %d, got %d", num, got)
}
