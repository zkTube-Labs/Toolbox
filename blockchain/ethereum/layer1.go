package eth

import (
	"context"
	"errors"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

func (E *Ethereum) AddressIsContract(address string) (isContract bool, err error) {
	c_address := common.HexToAddress(address)
	bytecode, err := E.Cli.CodeAt(context.Background(), c_address, nil) // nil is latest block
	if err != nil {
		return
	}
	isContract = len(bytecode) > 0
	return
}

func (E *Ethereum) GetBalance(address string) (balance *big.Int, err error) {
	return E.Cli.PendingBalanceAt(context.Background(), common.HexToAddress(address))
}

func (E *Ethereum) GetBlock(blockNumber *big.Int) (block *types.Block, err error) {
	return E.Cli.BlockByNumber(context.Background(), blockNumber)
}

func (E *Ethereum) TxETH(privKey, to string, toAmount *big.Float, gasLimit uint64, gasPrice *big.Int) (TxHash common.Hash, signedTx *types.Transaction, err error) {
	if gasPrice == nil {
		gasPrice, err = E.Cli.SuggestGasPrice(context.Background())
		if err != nil {
			return
		}
	}

	fromPrivkey, err := GetEcdsaByPrivatekey(privKey)
	if err != nil {
		return
	}
	addrHash, err := GetAddressByPrivateKey(fromPrivkey)
	if err != nil {
		return
	}
	balance, err := E.GetBalance(addrHash)
	if err != nil {
		return
	}
	amount := EthToWei(toAmount)
	if balance.Cmp(amount) < 0 {
		err = errors.New("insufficient balance")
		return
	}

	fromAddr := common.HexToAddress(addrHash)
	toAddr := common.HexToAddress(to)

	// get nonce
	nonce, err := E.Cli.PendingNonceAt(context.Background(), fromAddr)
	if err != nil {
		return
	}

	chainID, err := E.Cli.ChainID(context.Background())
	if err != nil {
		return
	}

	auth, err := bind.NewKeyedTransactorWithChainID(fromPrivkey, chainID)
	if err != nil {
		panic(err)
	}
	auth.Nonce = big.NewInt(int64(nonce))
	auth.Value = amount      // in wei
	auth.GasLimit = gasLimit // in units
	auth.GasPrice = gasPrice
	auth.From = fromAddr

	tx := types.NewTransaction(nonce, toAddr, amount, gasLimit, gasPrice, []byte{})
	signedTx, err = auth.Signer(fromAddr, tx)
	if err != nil {
		return
	}

	err = E.Cli.SendTransaction(context.Background(), signedTx)
	TxHash = tx.Hash()
	return
}
