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
	"encoding/hex"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	scom "github.com/ethereum/go-ethereum/contracts/native/cross_chain_manager/common"
	"github.com/ethereum/go-ethereum/contracts/native/utils"
	"github.com/ethereum/go-ethereum/core/state"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestUnpackMakeProofEvent(t *testing.T) {
	data := "0x000000000000000000000000000000000000000000000000000000000000006000000000000000000000000000000000000000000000000000000000000079de0000000000000000000000000000000000000000000000000000000000000220000000000000000000000000000000000000000000000000000000000000019a6638636261303633336164383763666533393033626636353935643437363035323066633761616363323032363862336431323337346438633866303665336633383462666130316638613761303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303461303065343365663531653131633332353730303030666361643465373065393432373735383965383062393762633435303135306331336164356231613537353239343764373964393336646137383333633766653035366562343530303634663334613332376463613834643934376437396439333664613738333363376665303536656234353030363466333461333237646361383834366436393665373462346633393430303030303030303030303030303030303030303030303030303030303030303030303030303030393430396634653438346434336233643662323039353766376531373630626565336336663632313836383830646530623662336137363430303030000000000000000000000000000000000000000000000000000000000000000000000000008635373437633035666632333666386431386262323162633032656363333839646566383533636165373236353731373536353733373434643030303030303030303030303030363333616438376366653339303362663635393564343736303532306663376161636332303236386233643132333734643863386630366533663338346266610000000000000000000000000000000000000000000000000000"
	blob, err := hexutil.Decode(data)
	assert.NoError(t, err)

	list, err := utils.UnpackEvent(*scom.ABI, scom.NOTIFY_MAKE_PROOF_EVENT, blob)
	assert.NoError(t, err)

	assert.Equal(t, len(list), 3)
	rawKey := list[2].(string)
	raw, err := hex.DecodeString(rawKey)
	assert.NoError(t, err)

	slot := state.Key2Slot(raw[common.AddressLength:])
	key := hexutil.Encode(slot[:])
	storageKeys := []string{key}
	t.Log(storageKeys)
}

func TestSimple(t *testing.T) {
	data := "0xf8cba0d0635e0ee3be359faf968a8d292ba39e50dc05dccf6d86c1213494a242b5247501f8a7a00000000000000000000000000000000000000000000000000000000000000007a03a3131469fcfa5e09ad7d8f455d2c45cbb707eff3ed97de7c3800971571b0dc8947d79d936da7833c7fe056eb450064f34a327dca84d947d79d936da7833c7fe056eb450064f34a327dca8846d696e74b4f39400000000000000000000000000000000000000009409f4e484d43b3d6b20957f7e1760bee3c6f62186880de0b6b3a7640000"
	dec, err := hexutil.Decode(data)
	assert.NoError(t, err)
	hash := crypto.Keccak256Hash(dec)
	t.Logf(hash.Hex())
}