package eth

import (
	"context"
	"errors"
	"math/big"

	token "github.com/zkTube-Labs/Toolbox/blockchain/ethereum/abi/erc20"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

var EmptyToken = common.HexToAddress("0x0000000000000000000000000000000000000000")

func (E *Ethereum) AddressIsContract(address string) (isContract bool, err error) {
	c_address := common.HexToAddress(address)
	bytecode, err := E.Cli.CodeAt(context.Background(), c_address, nil) // nil is latest block
	if err != nil {
		return
	}
	isContract = len(bytecode) > 0
	return
}

func (E *Ethereum) GetBalance(address common.Address, tokenAddr common.Address) (balance *big.Int, err error) {
	if tokenAddr != EmptyToken {
		instance, terr := token.NewErc20Caller(tokenAddr, E.Cli)
		if terr != nil {
			return nil, terr
		}
		balance, err = instance.BalanceOf(&bind.CallOpts{
			Pending: true,
		}, address)
	} else {
		balance, err = E.Cli.PendingBalanceAt(context.Background(), address)
	}
	return
}

func (E *Ethereum) GetBlock(blockNumber *big.Int) (block *types.Block, err error) {
	return E.Cli.BlockByNumber(context.Background(), blockNumber)
}

func (E *Ethereum) Transfer(privKey, to string, toAmount *big.Float, gasLimit uint64, gasPrice *big.Int, tokenHash string) (TxHash common.Hash, err error) {
	tokenAddr := common.HexToAddress(tokenHash)
	toAddr := common.HexToAddress(to)
	fromAddr, auth, err := E.createAuth(privKey, gasLimit, gasPrice)
	if err != nil {
		return
	}
	balance, err := E.GetBalance(fromAddr, tokenAddr)
	if err != nil {
		return
	}
	amount, err := E.getAmount(toAmount, tokenAddr)
	if err != nil {
		return
	}
	if balance.Cmp(amount) < 0 {
		err = errors.New("insufficient balance")
		return
	}
	if tokenAddr != EmptyToken {
		return E.transferERC20(auth, tokenAddr, toAddr, amount)
	} else {
		return E.transferETH(auth, toAddr, amount)
	}
}

func (E *Ethereum) getAmount(toAmount *big.Float, tokenAddr common.Address) (amount *big.Int, err error) {
	if tokenAddr != EmptyToken {
		instance, terr := token.NewErc20Caller(tokenAddr, E.Cli)
		if terr != nil {
			return nil, terr
		}
		decimals, terr := instance.Decimals(&bind.CallOpts{})
		if err != nil {
			return nil, terr
		}
		amount = AmountToToken(toAmount, int(decimals))
	} else {
		amount = EthToWei(toAmount)
	}

	return
}

func (E *Ethereum) createAuth(privKey string, gasLimit uint64, gasPrice *big.Int) (fromAddr common.Address, auth *bind.TransactOpts, err error) {
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
	fromAddr = common.HexToAddress(addrHash)

	nonce, err := E.Cli.PendingNonceAt(context.Background(), fromAddr)
	if err != nil {
		return
	}

	chainID, err := E.Cli.ChainID(context.Background())
	if err != nil {
		return
	}

	auth, err = bind.NewKeyedTransactorWithChainID(fromPrivkey, chainID)
	if err != nil {
		panic(err)
	}
	auth.Nonce = big.NewInt(int64(nonce))
	auth.GasLimit = gasLimit
	auth.GasPrice = gasPrice
	auth.From = fromAddr
	return
}

func (E *Ethereum) transferETH(auth *bind.TransactOpts, toAddr common.Address, amount *big.Int) (TxHash common.Hash, err error) {
	auth.Value = amount
	tx := types.NewTransaction(auth.Nonce.Uint64(), toAddr, auth.Value, auth.GasLimit, auth.GasPrice, []byte{})
	TxHash = tx.Hash()
	signedTx, err := auth.Signer(auth.From, tx)
	if err != nil {
		return
	}
	err = E.Cli.SendTransaction(context.Background(), signedTx)
	return
}

func (E *Ethereum) transferERC20(auth *bind.TransactOpts, tokenAddr, toAddr common.Address, amount *big.Int) (TxHash common.Hash, err error) {
	instance, err := token.NewErc20Transactor(tokenAddr, E.Cli)
	if err != nil {
		return
	}
	tx, err := instance.Transfer(auth, toAddr, amount)
	TxHash = tx.Hash()
	return
}
