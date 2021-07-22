package journal

import (
	"fmt"
	"math/big"
	"time"

	"github.com/dylenfu/zion-tool/config"
	"github.com/dylenfu/zion-tool/flag"
	"github.com/dylenfu/zion-tool/sdk"
	"github.com/dylenfu/zion-tool/utils/files"
	"github.com/dylenfu/zion-tool/utils/math"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/urfave/cli"
)

func getConfig(ctx *cli.Context) (*config.Config, error) {
	c := new(config.Config)
	cfgPath := flag.Flag2string(ctx, flag.ConfigPathFlag)
	if err := files.ReadJsonFile(cfgPath, c); err != nil {
		return nil, fmt.Errorf("read config json file, err: %v", err)
	}
	return c, nil
}

func getPeriodAndTxn(ctx *cli.Context) (int, int, error) {
	tLasting, err := flag.Flag2Duration(ctx, flag.PeriodFlag)
	if err != nil {
		return 0, 0, err
	}
	txn := int(flag.Flag2Uint64(ctx, flag.TxPerPeriod))
	period := int(tLasting / time.Second)
	return period, txn, nil
}

func getInstanceNumber(ctx *cli.Context) int {
	return int(flag.Flag2Uint64(ctx, flag.NumberFlag))
}

func generateMasterAccount(c *config.Config) (*sdk.Account, error) {
	masterPK, err := crypto.HexToECDSA(c.MasterNodeKey)
	if err != nil {
		return nil, fmt.Errorf("get main node key failed, err: %v", err)
	}
	return sdk.CustomNewAccount(c.ChainID, c.NodeList[0], masterPK)
}

func generateMultiTestingAccounts(c *config.Config, num int) ([]*sdk.Account, error) {
	accounts := make([]*sdk.Account, num)
	for idx := 0; idx < num; idx++ {
		url := c.NodeList[idx%len(c.NodeList)]
		acc, err := sdk.NewAccount(c.ChainID, url)
		if err != nil {
			return nil, err
		}
		accounts[idx] = acc
	}
	return accounts, nil
}

func prepareTestingAccountsBalance(master *sdk.Account, accounts []*sdk.Account, instanceNum, period, txn int) error {
	//logger := orlogger.New("prepare balance", "master", master.Address().Hex())

	amount := totalAmount(period, txn)
	gas := totalGas(period, txn)
	total := math.SafeAdd(amount, gas)
	for idx := 0; idx < len(accounts); idx++ {
		if _, err := master.Transfer(accounts[idx].Address(), total); err != nil {
			return err
		}
	}

	time.Sleep(5 * time.Second)

	for idx := 0; idx < instanceNum; idx++ {
		account := accounts[idx]
		balance, err := account.Balance(nil)
		if err != nil {
			return err
		}
		if balance.Cmp(total) < 0 {
			time.Sleep(10 * time.Second)
			//return fmt.Errorf("%s balance not engough", account.Address().Hex())
		} else {
			fmt.Println("deposit for account", "address", account.Address().Hex(), "balance", math.PrintUT(balance))
		}
	}

	return nil
}

func calculateTPS(master *sdk.Account, period int)  {
	//logger := orlogger.New("calculate tps", "period", period)

	startBlockNo, err := master.CurrentBlockNumber()
	if err != nil {
		panic(fmt.Sprintf("try to get start block number failed, err: %v", err))
	} else {
		fmt.Println("start from block", startBlockNo)
	}

	cnt := 0
	totalTx := uint(0)
	curBlockNum := startBlockNo
	startTime, endTime := uint64(0), uint64(0)
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
		endTime = header.Time

		retryTxCnt:
		txn, err := master.TxNum(header.Hash())
		if err != nil {
			time.Sleep(500 * time.Millisecond)
			goto retryTxCnt
		}
		totalTx += txn

		if endTime > startTime {
			tps := totalTx / uint((endTime - startTime))
			fmt.Println("calculate tps", "startBlock", startBlockNo, "endBlock", curBlockNum, "start time", startTime, "end time", endTime, "total tx", totalTx, "tps", tps)
		}

		curBlockNum += 1
		cnt += 1
	}
}

var (
	amountPerTx = big.NewInt(100000000000000) // 0.0001 eth
	extraGas    = math.SafeMul(big.NewInt(1), ETH1)
)

func totalTx(periodCnt, txPerPeriod int) *big.Int {
	return new(big.Int).SetUint64(uint64(periodCnt * txPerPeriod))
}

func totalAmount(periodCnt, txPerPeriod int) *big.Int {
	txn := totalTx(periodCnt, txPerPeriod)
	return math.SafeMul(amountPerTx, txn)
}

func defaultGas() *big.Int {
	return math.SafeMul(sdk.DefaultGasPrice(), sdk.DefaultGasLimit())
}

// totalGas = gasUsed * 2 + extraGas
func totalGas(periodCnt, txPerPeriod int) *big.Int {
	txn := totalTx(periodCnt, txPerPeriod)
	gasPerTx := defaultGas()
	total := math.SafeMul(gasPerTx, txn)
	extra := math.SafeAdd(total, extraGas)
	return math.SafeAdd(total, extra)
}
