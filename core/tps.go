package core

import (
	"fmt"
	"time"

	"github.com/dylenfu/zion-tool/pkg/math"
	"github.com/dylenfu/zion-tool/pkg/sdk"
	"github.com/ethereum/go-ethereum/common"
	"github.com/urfave/cli"
)

//var orlogger log.Logger

func Init() {
	math.Init(18)
	//orlogger = log.New("Handle TPS", ": ")
}

// HandleTPS try to test hotstuff tps, params nodeList represents multiple ethereum rpc url addresses,
// and num denote that this test will use multi account to send simple transaction
func HandleTPS(ctx *cli.Context) error {
	fmt.Println("start to handle tps", "start", true)

	// load config instance
	c, err := getConfig(ctx)
	if err != nil {
		return err
	}

	// load and try to increase gas price
	setGasPriceIncr(ctx)

	// load period and tx number per period
	fmt.Println("try to get period and txn...")
	period, txn, err := getPeriodAndTxn(ctx)
	if err != nil {
		return err
	}
	fmt.Println("get period and txn", "period", period, "txn", txn)

	// generate master account
	fmt.Println("try to generate master account...")
	master, err := generateMasterAccount(c)
	if err != nil {
		return err
	}
	fmt.Println("generate master account", "period", period, "txn", txn)

	// create account
	fmt.Println("try to generate multi test accounts...")
	instanceNo := getInstanceNumber(ctx)
	accounts, err := generateMultiTestingAccounts(c, instanceNo)
	if err != nil {
		return err
	}
	fmt.Println("generated multi test accounts")

	// prepare balance
	fmt.Println("try to prepare test accounts balance...")
	if err := prepareTestingAccountsBalance(master, accounts, instanceNo, period, txn); err != nil {
		return err
	}
	fmt.Println("prepare test accounts balance success")

	// send transactions continuously
	to := master.Address()
	for _, acc := range accounts {
		go sendTransfer(acc, to, txn)
	}

	// calculate and print tps
	calculateTPS(master, period)
	return nil
}

func sendTransfer(acc *sdk.Account, to common.Address, txn int) {
	ticker := time.NewTicker(1 * time.Second)

	for range ticker.C {
		for i := 0; i < txn; i++ {
			if _, err := acc.Transfer(to, amountPerTx); err != nil {
				fmt.Println("transfer failed", "err", err)
			} else {
				//fmt.Println("transfer success", "hash", hash)
			}
		}
	}
}
