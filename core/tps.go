package core

import (
	"fmt"
	"math/big"
	"time"

	"github.com/dylenfu/zion-tool/config"
	"github.com/dylenfu/zion-tool/pkg/encode"
	"github.com/dylenfu/zion-tool/pkg/log"
	"github.com/dylenfu/zion-tool/pkg/math"
	"github.com/dylenfu/zion-tool/pkg/sdk"
	"github.com/ethereum/go-ethereum/common"
)

var (
	ETH1 = math.Pow10toBigInt(18)
)

// TPS try to test hotstuff tps, params nodeList represents multiple ethereum rpc url addresses,
// and num denote that this test will use multi account to send simple transaction
func TPS() bool {
	log.Info("start to handle tps")

	var params struct {
		AccountsNum int
		LastTime    encode.Duration
		TxPerSecond int
		InstanceNum int
		IncreaseGas int
	}

	// load config instance
	if err := config.LoadParams("test_tps.json", &params); err != nil {
		log.Error(err)
		return false
	}

	// generate master account
	log.Info("try to generate master account...")
	master, err := masterAccount()
	if err != nil {
		log.Errorf("load master account failed, err: %v", err)
		return false
	}
	log.Split("generate master account success!")

	// create account
	log.Info("try to generate multi test accounts...")
	nodes := config.Conf.Nodes[:params.InstanceNum]
	accounts, err := generateMultiTestingAccounts(nodes, params.AccountsNum)
	if err != nil {
		log.Errorf("generate multi testing accounts failed, err: %v", err)
		return false
	}
	log.Split("generated multi test accounts success!")

	// prepare balance
	log.Info("try to prepare test accounts balance...")
	period := int(time.Duration(params.LastTime) / time.Second)
	if err := prepareTestingAccountsBalance(master, accounts, period, params.TxPerSecond); err != nil {
		log.Errorf("prepare testing accounts balance failed, err: %v", err)
		return false
	}
	log.Split("prepare test accounts balance success!")

	// send transactions continuously
	to := master.Address()
	for _, acc := range accounts {
		go sendTransfer(acc, to, params.TxPerSecond)
	}

	// calculate and print tps
	calculateTPS(master, period)
	return true
}

func sendTransfer(acc *sdk.Account, to common.Address, txn int) {
	var duration int = 5
	txsNum := txn * duration
	ticker := time.NewTicker(time.Duration(duration) * time.Second)

	for range ticker.C {
		for i := 0; i < txsNum; i++ {
			if _, err := acc.Transfer(to, amountPerTx); err != nil {
				log.Errorf("transfer failed", "err", err)
			}
		}
	}
}

func generateMultiTestingAccounts(nodes []*config.Node, num int) ([]*sdk.Account, error) {
	if nodes == nil || len(nodes) == 0 {
		return nil, fmt.Errorf("invalid nodes")
	}

	chainID := config.Conf.ChainID
	accounts := make([]*sdk.Account, num)
	nodesLen := len(nodes)
	for idx := 0; idx < num; idx++ {
		url := nodes[idx%nodesLen].Url
		acc, err := sdk.NewAccount(chainID, url)
		if err != nil {
			return nil, err
		}
		accounts[idx] = acc
	}
	return accounts, nil
}

func prepareTestingAccountsBalance(master *sdk.Account, accounts []*sdk.Account, period, txn int) error {
	amount := totalAmount(period, txn)
	gas := totalGas(period, txn)
	total := math.SafeAdd(amount, gas)
	balanceMap := make(map[string]*sdk.Account)
	for idx := 0; idx < len(accounts); idx++ {
		addr := accounts[idx].Address()
		if tx, err := master.Transfer(addr, total); err != nil {
			return err
		} else {
			log.Info("master transfer", total, "to", addr.Hex(), "hash", tx.Hex())
		}
		balanceMap[addr.Hex()] = accounts[idx]
	}

	time.Sleep(5 * time.Second)

retry:
	for addr, account := range balanceMap {
		balance, err := account.Balance(nil)
		if err != nil {
			return err
		}
		if balance.Cmp(total) >= 0 {
			log.Info("deposit for account", "address", account.Address().Hex(), "balance", math.PrintUT(balance))
			delete(balanceMap, addr)
		}
	}

	if len(balanceMap) > 0 {
		time.Sleep(5 * time.Second)
		log.Infof("there are %d account need to preparing", len(balanceMap))
		goto retry
	}

	return nil
}

func calculateTPS(master *sdk.Account, period int) {
	startBlockNo, err := master.CurrentBlockNumber()
	if err != nil {
		panic(fmt.Sprintf("try to get start block number failed, err: %v", err))
	} else {
		log.Info("start from block", startBlockNo)
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
			log.Info("calculate tps", "startBlock", startBlockNo, "endBlock", curBlockNum, "start time", startTime, "end time", endTime, "total tx", totalTx, "tps", tps)
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
