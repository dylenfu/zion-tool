package journal

import (
	"fmt"
	"time"

	"github.com/urfave/cli"
)

func PolyChainListen(ctx *cli.Context) error {
	fmt.Println("start to listen", "start", true)
	//获取config
	c, err := getConfig(ctx)
	if err != nil {
		return err
	}
	master, err := generateMasterAccount(c)
	if err != nil {
		return err
	}
	startBlockNo, err := master.CurrentBlockNumber()
	if err != nil {
		panic(fmt.Sprintf("try to get start block number failed, err: %v", err))
	} else {
		fmt.Println("start from block", startBlockNo)
	}
	
	period, txn, err := getPeriodAndTxn(ctx)
	if err != nil {
		return err
	}
	fmt.Println("get period and txn", "period", period, "txn", txn)

	cnt := 0
	totalTx := uint(0)
	curBlockNum := startBlockNo
	startTime, endTime, preTime := uint64(0), uint64(0), uint64(0)

	for cnt < period {
	retryHeader:
		header, err := master.BlockHeaderByNumber(curBlockNum)
		if err != nil {
			time.Sleep(500 * time.Millisecond)
			goto retryHeader
		}
		if curBlockNum == startBlockNo {
			startTime = header.Time
		}

		preTime = endTime
		endTime = header.Time

	retryPendingTX:
		pendingTxNum, err := master.PendingTransactionNum()
		if err != nil {
			time.Sleep(500 * time.Millisecond)
			goto retryPendingTX
		}

	retryTxCnt:
		txn, err := master.TxNum(header.Hash())
		if err != nil {
			time.Sleep(500 * time.Millisecond)
			goto retryTxCnt
		}
		totalTx = txn

		if endTime > startTime {
			tps := totalTx / uint((endTime - preTime))
			fmt.Println("calculate tps", "startBlock", startBlockNo, "endBlock", curBlockNum, "pre time", preTime, "end time", endTime, "pendingTx NUM", pendingTxNum, "total tx", totalTx, "tps", tps)
		}

		curBlockNum += 1
		cnt += 1
	}
	return nil
}
