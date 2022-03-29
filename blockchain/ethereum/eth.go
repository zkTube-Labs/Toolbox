package eth

import (
	"crypto/ecdsa"
	"errors"
	"math"
	"math/big"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

const (
	Mainnet Network = "mainnet"
	Rinkeby Network = "rinkeby"
	Ropsten Network = "ropsten"

	Ifr = "infura"
)

type Network string

var ETH *Ethereum

type Ethereum struct {
	HotWallet string
	ColdWallt string
	Infura    *infura
	Cli       *ethclient.Client
	WSCli     *ethclient.Client
}

func InitEth(Hot, Cold string) (err error) {
	ETH = &Ethereum{
		HotWallet: Hot,
		ColdWallt: Cold,
		Infura:    NewInfura(),
	}
	ETH.Cli, err = ETH.Infura.GetClient()
	return
}

func NewEth() *Ethereum {
	return ETH
}

func NewPrivateKey() (privateKey *ecdsa.PrivateKey, err error) {
	privateKey, err = crypto.GenerateKey()
	return
}

func GetPriaveKeyByEcdsa(privateKey *ecdsa.PrivateKey) (privateKeyStr string) {
	privateKeyBytes := crypto.FromECDSA(privateKey)
	privateKeyStr = hexutil.Encode(privateKeyBytes)
	return
}

func GetEcdsaByPrivatekey(privateKeyStr string) (privateKey *ecdsa.PrivateKey, err error) {
	return crypto.HexToECDSA(privateKeyStr)
}

func GetAddressByPrivateKey(privateKey *ecdsa.PrivateKey) (address string, err error) {
	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		err = errors.New("cannot assert type: publicKey is not of type *ecdsa.PublicKey")
		return
	}
	address = crypto.PubkeyToAddress(*publicKeyECDSA).Hex()
	return
}

func EthToWei(eth *big.Float) *big.Int {
	return AmountToToken(eth, 18)
}

func WeiToEth(wei *big.Int) *big.Float {
	return TokenToAmount(wei, 18)
}

func AmountToToken(amount *big.Float, decimals int) (i *big.Int) {
	i = new(big.Int)
	multiple := new(big.Float).SetFloat64(math.Pow10(decimals))
	amount.Mul(amount, multiple).Int(i)
	return
}

func TokenToAmount(token *big.Int, decimals int) *big.Float {
	divisor := new(big.Float).SetInt(token)
	dividend := new(big.Float).SetFloat64(math.Pow10(decimals))
	return new(big.Float).Quo(divisor, dividend)
}
