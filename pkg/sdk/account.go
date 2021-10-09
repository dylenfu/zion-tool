package sdk

import (
	"context"
	"crypto/ecdsa"
	"math/big"
	"sync"

	"github.com/ethereum/go-ethereum/rpc"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/contracts/native/governance"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

var (
	gasLimit = uint64(210000)
	gasPrice = new(big.Int).SetUint64(1000000000)

	EmptyHash = common.Hash{}
	goverABI  = governance.GetABI()
)

type Account struct {
	signer    types.EIP155Signer
	pk        *ecdsa.PrivateKey
	address   common.Address
	url       string
	client    *ethclient.Client
	rpcClient *rpc.Client

	nonce   uint64
	nonceMu *sync.RWMutex
}

func NewAccount(chainID uint64, url string) (*Account, error) {
	pk, err := crypto.GenerateKey()
	if err != nil {
		return nil, err
	}

	return CustomNewAccount(chainID, url, pk)
}

func CustomNewAccount(chainID uint64, url string, pk *ecdsa.PrivateKey) (*Account, error) {
	address := crypto.PubkeyToAddress(pk.PublicKey)
	signer := types.NewEIP155Signer(new(big.Int).SetUint64(chainID))
	rpcclient, err := rpc.Dial(url)
	if err != nil {
		return nil, err
	}
	client := ethclient.NewClient(rpcclient)

	curNonce, err := client.NonceAt(context.Background(), address, nil)
	if err != nil {
		return nil, err
	}

	acc := &Account{
		signer:    signer,
		pk:        pk,
		address:   address,
		url:       url,
		client:    client,
		rpcClient: rpcclient,
		nonce:     curNonce,
		nonceMu:   new(sync.RWMutex),
	}
	return acc, nil
}

func (c *Account) Address() common.Address {
	return c.address
}

func (c *Account) Url() string {
	return c.url
}

func (c *Account) Balance(blockNum *big.Int) (*big.Int, error) {
	return c.client.BalanceAt(context.Background(), c.address, blockNum)
}

func (c *Account) BalanceOf(addr common.Address, blockNum *big.Int) (*big.Int, error) {
	return c.client.BalanceAt(context.Background(), addr, blockNum)
}

func (c *Account) Transfer(to common.Address, amount *big.Int) (common.Hash, error) {
	signedTx, err := c.NewSignedTx(to, amount, nil)
	if err != nil {
		return EmptyHash, err
	}
	if err := c.SendTx(signedTx); err != nil {
		return EmptyHash, err
	}
	return signedTx.Hash(), nil
}

func (c *Account) Nonce() uint64 {
	c.nonceMu.RLock()
	defer c.nonceMu.RUnlock()
	return c.nonce
}

func (c *Account) NewUnsignedTx(to common.Address, amount *big.Int, data []byte) (*types.Transaction, error) {
	nonce := c.Nonce()
	gasLimit := DefaultGasLimit().Uint64()
	gasPrice := big.NewInt(2000000000)
	//gasPrice, err := c.client.SuggestGasPrice(context.Background())
	//if err != nil {
	//	return nil, err
	//}
	return types.NewTx(&types.LegacyTx{
		Nonce:    nonce,
		To:       &to,
		Value:    amount,
		Gas:      gasLimit,
		GasPrice: gasPrice,
		Data:     data,
	}), nil
	// return types.NewTransaction(nonce, to, amount, gas, price, data)
}

func (c *Account) NewSignedTx(to common.Address, amount *big.Int, data []byte) (*types.Transaction, error) {
	unsignedTx, err := c.NewUnsignedTx(to, amount, data)
	if err != nil {
		return nil, err
	}
	return types.SignTx(unsignedTx, c.signer, c.pk)
}

func (c *Account) SendTx(signedTx *types.Transaction) error {
	defer func() {
		c.nonceMu.Lock()
		c.nonce += 1
		c.nonceMu.Unlock()
	}()

	return c.client.SendTransaction(context.Background(), signedTx)
}

func (c *Account) CurrentBlockNumber() (uint64, error) {
	return c.client.BlockNumber(context.Background())
}

func (c *Account) BlockHeaderByNumber(blockNumber uint64) (*types.Header, error) {
	return c.client.HeaderByNumber(context.Background(), new(big.Int).SetUint64(blockNumber))
}

func (c *Account) TxNum(blockHash common.Hash) (uint, error) {
	return c.client.TransactionCount(context.Background(), blockHash)
}

func (c *Account) CallContract(caller, contractAddr common.Address, payload []byte, blockNum string) ([]byte, error) {
	arg := ethereum.CallMsg{
		From: caller,
		To:   &contractAddr,
		Data: payload,
	}

	// todo: block number
	return c.client.CallContract(context.Background(), arg, nil)
}

func AddGasPrice(inc uint64) {
	added := new(big.Int).SetUint64(inc)
	gasPrice = new(big.Int).Add(gasPrice, added)
}

func DefaultGasPrice() *big.Int {
	return gasPrice
}

func DefaultGasLimit() *big.Int {
	return new(big.Int).SetUint64(gasLimit)
}
