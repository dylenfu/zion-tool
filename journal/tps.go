package journal

import (
	"sync"
	"time"

	"github.com/urfave/cli"
)

// HandleTPS try to test hotstuff tps, params nodeList represents multiple ethereum rpc url addresses,
// and num denote that this test will use multi account to send simple transaction
func HandleTPS(ctx *cli.Context) error {
	// load config instance
	c, err := getConfig(ctx)
	if err != nil {
		return err
	}

	// load period and tx number per period
	period, txn, err := getPeriodAndTxn(ctx)
	if err != nil {
		return err
	}

	// generate master account
	master, err := generateMasterAccount(c)
	if err != nil {
		return err
	}

	// create account
	instanceNo := getInstanceNumber(ctx)
	accounts, err := generateMultiTestingAccounts(c, instanceNo)
	if err != nil {
		return err
	}

	// prepare balance
	if err := prepareTestingAccountsBalance(master, accounts, instanceNo, period, txn); err != nil {
		return err
	}

	// send transactions continuously
	to := master.Address()
	wg := new(sync.WaitGroup)
	for _, acc := range accounts {
		go func() {
			periodCnt := 0
			ticker := time.NewTicker(1 * time.Second)

			sendMultiTransfer := func() {
				for i := 0; i < txn; i++ {
					acc.Transfer(to, amountPerTx)
				}
			}

			for range ticker.C {
				if periodCnt += 1; periodCnt < period {
					sendMultiTransfer()
				}
			}
			wg.Add(1)
		}()
	}
	wg.Done()

	return nil
}
