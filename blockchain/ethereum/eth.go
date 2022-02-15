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
	Mainnet = "mainnet"
	Rinkeby = "rinkeby"
	Ropsten = "ropsten"

	Ifr = "infura"
)

var ETH *Ethereum

type Ethereum struct {
	HotWallet string
	ColdWallt string
	Infura    *infura
	Cli       *ethclient.Client
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

func WeiToEth(wei *big.Int) *big.Float {
	divisor := new(big.Float).SetInt(wei)
	dividend := new(big.Float).SetFloat64(math.Pow10(18))
	return new(big.Float).Quo(divisor, dividend)
}

func EthToWei(eth *big.Float) *big.Int {
	multiple := new(big.Float).SetFloat64(math.Pow10(18))
	ui, _ := new(big.Float).Mul(eth, multiple).Uint64()
	return new(big.Int).SetUint64(ui)
}
