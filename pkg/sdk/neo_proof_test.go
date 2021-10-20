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

	. "github.com/dylenfu/zion-tool/pkg/go_abi/neo_proof"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/contracts/native/utils"
	"github.com/stretchr/testify/assert"
)

func TestDecodeNeoProofMethodSetInput(t *testing.T) {
	input := "0xe942b51600000000000000000000000000000000000000000000000000000000000000400000000000000000000000000000000000000000000000000000000000000080000000000000000000000000000000000000000000000000000000000000000668656c6c6f3800000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000006776f726c64380000000000000000000000000000000000000000000000000000"
	enc, err := hexutil.Decode(input)
	assert.NoError(t, err)
	var data struct {
		Key   string
		Value string
	}
	if err := utils.UnpackMethod(neoProofABI, MethodSet, &data, enc); err != nil {
		t.Fatal(err)
	}
	t.Logf("key %s, value %s", data.Key, data.Value)
}
