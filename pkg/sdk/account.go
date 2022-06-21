package sdk

import (
	"context"
	"crypto/ecdsa"
	"encoding/json"
	"fmt"
	"math/big"
	"sync"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/ethereum/go-ethereum/rpc"
)

var (
	gasLimit = uint64(210000)
	gasPrice = new(big.Int).SetUint64(1000000000)

	EmptyHash = common.Hash{}
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
	rpcclient, err := rpc.Dial(url)
	if err != nil {
		return nil, err
	}
	client := ethclient.NewClient(rpcclient)

	acc := &Account{
		pk:        pk,
		url:       url,
		client:    client,
		rpcClient: rpcclient,
	}

	if pk != nil {
		address := crypto.PubkeyToAddress(pk.PublicKey)
		signer := types.NewEIP155Signer(new(big.Int).SetUint64(chainID))
		curNonce, err := client.NonceAt(context.Background(), address, nil)
		if err != nil {
			return nil, err
		}
		acc.signer = signer
		acc.address = address
		acc.nonce = curNonce
		acc.nonceMu = new(sync.RWMutex)
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
	if err := c.WaitTransaction(signedTx.Hash()); err != nil {
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
	gasPrice, err := c.client.SuggestGasPrice(context.Background())
	if err != nil {
		return nil, err
	}

	callMsg := ethereum.CallMsg{
		From:     c.Address(),
		To:       &to,
		Gas:      0,
		GasPrice: gasPrice,
		Value:    amount,
		Data:     data,
	}
	gasLimit, err := c.client.EstimateGas(context.Background(), callMsg)
	if err != nil {
		return nil, fmt.Errorf("estimate gas limit error: %s", err.Error())
	}

	return types.NewTx(&types.LegacyTx{
		Nonce:    nonce,
		To:       &to,
		Value:    amount,
		Gas:      gasLimit,
		GasPrice: gasPrice,
		Data:     data,
	}), nil
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

func (c *Account) GetAccountAndStorageProof(contract common.Address, storageKeys []string, blockNum *big.Int) ([]byte, []byte, error) {
	proof, err := c.client.ProofAt(context.Background(), contract, storageKeys, blockNum)
	if err != nil {
		return nil, nil, err
	}
	if len(proof.StorageProof) < 1 {
		return nil, nil, fmt.Errorf("storage length invalid")
	}

	accountPrf, err := rlpEncodeStringList(proof.AccountProof)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to rlp account proof, err: %v", err)
	}
	storageProof, err := rlpEncodeStringList(proof.StorageProof[0].Proof)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to rlp storage proof, err: %v", err)
	}

	return accountPrf, storageProof, nil
}

func (c *Account) StorageAt(contract common.Address, storageKey common.Hash, blockNum *big.Int) ([]byte, error) {
	return c.client.StorageAt(context.Background(), contract, storageKey, blockNum)
}

func (c *Account) GetProof(contract common.Address, storageKeys []string, blockNum *big.Int) ([]byte, error) {
	proof, err := c.client.ProofAt(context.Background(), contract, storageKeys, blockNum)
	if err != nil {
		return nil, err
	}
	return json.Marshal(proof)
}

func rlpEncodeStringList(raw []string) ([]byte, error) {
	var rawBytes []byte
	for i := 0; i < len(raw); i++ {
		rawBytes = append(rawBytes, common.Hex2Bytes(raw[i][2:])...)
	}
	return rlp.EncodeToBytes(rawBytes)
}

func (c *Account) CallContract(caller, contractAddr common.Address, payload []byte, blockNum *big.Int) ([]byte, error) {
	arg := ethereum.CallMsg{
		From: caller,
		To:   &contractAddr,
		Data: payload,
	}

	return c.client.CallContract(context.Background(), arg, blockNum)
}

func (c *Account) signAndSendTx(payload []byte, contract common.Address) (common.Hash, error) {
	return c.signAndSendTxWithValue(payload, big.NewInt(0), contract)
}

func (c *Account) signAndSendTxWithValue(payload []byte, amount *big.Int, contract common.Address) (common.Hash, error) {
	hash := common.EmptyHash
	tx, err := c.NewSignedTx(contract, amount, payload)
	if tx != nil {
		hash = tx.Hash()
	}
	if err != nil {
		return hash, fmt.Errorf("sign tx failed, err: %v", err)
	}

	if err := c.SendTx(tx); err != nil {
		return hash, err
	}
	if err := c.WaitTransaction(tx.Hash()); err != nil {
		return hash, err
	}
	return hash, nil
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
