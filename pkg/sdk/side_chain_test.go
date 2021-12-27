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
	"testing"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	lp "github.com/ethereum/go-ethereum/contracts/native/go_abi/main_chain_lock_proxy_abi"
	"github.com/ethereum/go-ethereum/contracts/native/utils"
	"github.com/ethereum/go-ethereum/ethclient"
)

func TestFilterCrossChainEvent(t *testing.T) {
	url := "http://101.32.99.70:22000"
	height := uint64(3557568)
	address := utils.LockProxyContractAddress
	//blockHash := common.HexToHash("0xb4e5ac62c9e659b1d0fc8091f31140522b27d2e0aacea0ef0c501da61dc35c1c")

	cli, err := ethclient.Dial(url)
	if err != nil {
		t.Fatal(err)
	}
	eccm, err := lp.NewIMainChainLockProxy(address, cli)
	if err != nil {
		t.Fatal(err)
	}

	iterator, err := eccm.FilterCrossChainEvent(&bind.FilterOpts{
		Start:   height,
		End:     &height,
		Context: context.Background(),
	}, []common.Address{})
	if err != nil {
		t.Fatal(err)
	}

	for iterator.Next() {
		data := iterator.Event
		t.Logf("data is %v", data)
	}
	//ab, err := abi.JSON(strings.NewReader(lp.IMainChainLockProxyABI))
	//if err != nil {
	//	t.Fatal(err)
	//}
	//topic := ab.Events["CrossChainEvent"].ID
	//logs, err := cli.FilterLogs(context.Background(), ethereum.FilterQuery{
	//	BlockHash: &blockHash,
	//	Addresses: []common.Address{address},
	//	Topics:    [][]common.Hash{[]common.Hash{topic}},
	//})
	//if err != nil {
	//	t.Fatal(err)
	//}
	//t.Logf("real topic %s", topic.Hex())
	//t.Log(logs)
}
