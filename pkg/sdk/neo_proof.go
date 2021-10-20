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
	"fmt"
	"math/big"
	"strings"

	. "github.com/dylenfu/zion-tool/pkg/go_abi/neo_proof"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/contracts/native/utils"
)

var (
	neoProofAddr = common.HexToAddress("0xF59b9838a73CBCDebBF355f9dBD554435b672432")
	neoProofABI  = GetNeoProofABI()
)

func GetNeoProofABI() *abi.ABI {
	ab, err := abi.JSON(strings.NewReader(ProofABI))
	if err != nil {
		panic(err)
	}
	return &ab
}

func (c *Account) NeoProof(key, value string) (common.Hash, error) {
	payload, err := utils.PackMethod(neoProofABI, MethodSet, key, value)
	if err != nil {
		return common.EmptyHash, err
	}

	return c.sendNeoProofTx(payload)
}

func (c *Account) sendNeoProofTx(payload []byte) (common.Hash, error) {
	hash := common.EmptyHash
	tx, err := c.NewSignedTx(neoProofAddr, big.NewInt(0), payload)
	if tx != nil {
		hash = tx.Hash()
	}
	if err != nil {
		return hash, fmt.Errorf("sign tx failed, err: %v", err)
	}

	if err := c.SendTx(tx); err != nil {
		return hash, err
	}
	if err := c.WaitTransaction(tx.Hash()); err != nil {
		return hash, err
	}
	return hash, nil
}
