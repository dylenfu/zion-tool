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
	"github.com/dylenfu/zion-tool/pkg/frame"
	"github.com/dylenfu/zion-tool/pkg/math"
)

func Endpoint() {
	math.Init(18)

	frame.Tool.RegMethod("demo", Demo)
	frame.Tool.RegMethod("neo-proof", NeoProof)

	// main chain normal operation
	frame.Tool.RegMethod("transfer", Transfer)
	frame.Tool.RegMethod("tps", TPS)

	// main chain change bookeepers
	frame.Tool.RegMethod("epoch", Epoch)
	frame.Tool.RegMethod("history", EpochHistory)

	// main chain cross chain operation
	frame.Tool.RegMethod("reg-side-chain", RegisterSideChain)
	frame.Tool.RegMethod("approve-side-chain", ApproveSideChain)

	// sync side chain genesis header to main chain
	frame.Tool.RegMethod("sync-genesis-header", SyncGenesisHeader)

	// fetch main chain epoch info and sync to side chain eccm contract
	frame.Tool.RegMethod("epoch-proof", FetchEpochProof)

	// main chain mint
	frame.Tool.RegMethod("mint", Mint)
	frame.Tool.RegMethod("burn", Burn)
}

func Demo() bool {
	return true
}
